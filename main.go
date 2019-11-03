package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	dataDir     = "data"
	playlistDir = "data/playlists"
)

var (
	playlistItems *template.Template
	mediaItems    *template.Template

	webPathRE = regexp.MustCompile(`\.\.`)
)

func convertToWebPath(input string) (output string) {
	return webPathRE.ReplaceAllString(input, dataDir)
}

// IndexData is the data that index.tmpl.html takes.
type IndexData struct {
	Title         string
	PlaylistItems *PlaylistData
	MediaItems    *MediaData
}

func main() {
	log.SetFlags(log.Lshortfile)

	index, err := template.ParseFiles("index.tmpl.html")
	if err != nil {
		log.Fatalf("Parse index.tmpl.html: %v", err)
	} else if playlistItems = index.Lookup("Playlists"); playlistItems == nil {
		log.Fatalf("Could not lookup %q.", "Playlists")
	} else if mediaItems = index.Lookup("Media"); mediaItems == nil {
		log.Fatalf("Could not lookup %q.", "Media")
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		var (
			playlist = c.Query("p")
			media    = c.Query("m")
			_        = c.Query("t")
			err      error
		)
		m := &IndexData{
			Title: "GOPF: The GNU Online Player Framework",
		}
		if m.PlaylistItems, err = buildPlaylistData(playlist); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else if m.MediaItems, err = buildMediaData(playlist, media); err != nil {
			log.Printf("buildMediaData(): %v", err)
		}

		var buf bytes.Buffer
		if err := index.Execute(&buf, m); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Data(http.StatusOK, "text/html", buf.Bytes())
	})

	r.GET("/list", func(c *gin.Context) {
		if op := c.Query("op"); op != "" {
			switch op {
			case "ls":
				if dir := c.Query("dir"); dir != "" {
					if matches, err := filepath.Glob(filepath.Join(dir, "*")); err != nil {
						c.AbortWithError(http.StatusInternalServerError, err)
					} else {
						c.String(http.StatusOK, strings.Join(matches, "\n"))
					}
				} else {
					c.JSON(http.StatusBadRequest, map[string]string{"error": "Expected \"dir\" key."})
				}
				return

			case "playlist":
				generatePlaylists(c.Query("playlist"))

			default:
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("op: %s", op))
			}
		} else if playlist := c.Query("playlist"); playlist != "" {
			playlist = filepath.Join(playlistDir, playlist)
			if f, err := os.Open(playlist + ".json"); err == nil {
				defer f.Close()

				c.Writer.Header().Add("Content-Type", "application/json")
				if _, err = io.Copy(c.Writer, f); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				}
			} else if err != nil && !os.IsNotExist(err) {
				c.AbortWithError(http.StatusInternalServerError, err)
			} else if f, err := os.Open(playlist); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				if os.IsNotExist(err) {
					c.AbortWithError(http.StatusNotFound, fmt.Errorf("no such playlist: %v", playlist))
				} else {
					c.AbortWithError(http.StatusInternalServerError, err)
				}
				return
			} else {
				defer f.Close()

				c.Writer.Header().Add("Content-Type", "text/plain")
				if _, err := io.Copy(c.Writer, f); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
				}
			}
		} else {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Missing query param.",
			})
		}
	})

	r.NoRoute(func(c *gin.Context) {
		fname := filepath.Join(".", c.Request.URL.Path)
		f, err := os.Open(fname)
		if err != nil {
			if os.IsNotExist(err) {
				c.AbortWithError(http.StatusNotFound, err)
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		defer f.Close()

		ct := "text/plain"
		switch filepath.Ext(f.Name()) {
		case ".js":
			ct = "application/json"
		case ".css":
			ct = "text/css"
		}
		c.DataFromReader(http.StatusOK, -1, ct, f, nil)
	})

	r.Run(":8000")
}

// PlaylistData is used by the playlistItems template.
type PlaylistData struct {
	Playlist string
	Fnames   []string
}

func buildPlaylistData(playlist string) (p *PlaylistData, err error) {
	p = &PlaylistData{
		Playlist: playlist,
		Fnames:   make([]string, 0, 20),
	}
	err = filepath.Walk(playlistDir, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.Name()[0] == '.' || info.Mode()&0111 != 0 || info.Name()[len(info.Name())-1] == '~' {
			return nil
		}

		p.Fnames = append(p.Fnames, filepath.Base(info.Name()))
		return nil
	}))
	return
}

// Generates playlist from playlist files.  Playlist files are simple text files with each line
// containing the relative path to a song.
//
//   playlist: String name of playlist to start selected.
//
// Returns the HTML for the playlist list.
func generatePlaylists(playlist string) (string, error) {
	dat, err := buildPlaylistData(playlist)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := playlistItems.Execute(&buf, dat); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// MediaData is the structure needed by the mediaItems template.
type MediaData struct {
	Media   string
	Entries []MediaEntry
}

// MediaEntry is an entry for MediaData.
type MediaEntry struct {
	Name    string
	WebPath string
}

func buildMediaData(playlist, media string) (m *MediaData, err error) {
	m = &MediaData{
		Media: media,
	}

	b, err := ioutil.ReadFile(filepath.Join(playlistDir, playlist))
	if err != nil {
		return m, fmt.Errorf("reading file %q: %v", playlist, err)
	}

	fnames := strings.Split(string(b), "\n")
	for _, fname := range fnames {
		if fname = strings.TrimSpace(fname); len(fname) == 0 {
			continue
		}

		m.Entries = append(m.Entries, MediaEntry{
			Name:    filepath.Base(fname),
			WebPath: convertToWebPath(fname),
		})
	}

	return
}

// Generates the list of media in the media list.  This list could contain
// songs or videos.
//
//   playlist: The name of the playlist.
//   media: The name of the selected media (if any).
//
// Returns the HTML for the contents of the playlist with the passed name.
func generateMedia(playlist, media string) (string, error) {
	dat, err := buildMediaData(playlist, media)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err = mediaItems.Execute(&buf, dat); err != nil {
		return "", err
	}

	return buf.String(), nil
}
