package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

var (
	playlistItems *template.Template
	mediaItems    *template.Template

	dataDirs     []string
	playlistDirs []string

	webPathRE = regexp.MustCompile(`\.\.`)
)

// IndexData is the data that index.tmpl.html takes.
type IndexData struct {
	Title         string
	PlaylistItems *PlaylistData
	MediaItems    *MediaData
}

func main() {
	log.SetFlags(log.Lshortfile)

	var (
		dirs   = flag.String("data-dirs", "audio,video", "Comma-separated set of data directories.")
		prefix = flag.String("uri-prefix", os.Getenv("GOPF_URI_PREFIX"), "Set absolute URI prefix.")
		port   = flag.Int("port", 8000, "Port to listen on.")
	)
	flag.Parse()

	dataDirs = strings.Split(*dirs, ",")
	for i, dd := range dataDirs {
		dataDirs[i] = strings.TrimSpace(dd)
		playlistDirs = append(playlistDirs, filepath.Join(dataDirs[i], "playlists"))
	}

	makeURI := func(uri string) string {
		if *prefix == "" {
			return uri
		}
		return path.Join(*prefix, uri)
	}

	index, err := template.ParseFiles("index.tmpl.html")
	if err != nil {
		log.Fatalf("Parse index.tmpl.html: %v", err)
	} else if playlistItems = index.Lookup("Playlists"); playlistItems == nil {
		log.Fatalf("Could not lookup %q.", "Playlists")
	} else if mediaItems = index.Lookup("Media"); mediaItems == nil {
		log.Fatalf("Could not lookup %q.", "Media")
	}

	r := gin.Default()

	r.GET(makeURI("/"), func(c *gin.Context) {
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

	r.GET(makeURI("list"), func(c *gin.Context) {
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
		}

		playlist := c.Query("playlist")
		if playlist == "" {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Missing query param.",
			})
			return
		}

		var (
			b  []byte
			i  int
			pd string
		)
		for i, pd = range playlistDirs {
			if b, err = ioutil.ReadFile(filepath.Join(pd, playlist)); err != nil && !os.IsNotExist(err) {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			} else if err == nil {
				break
			}
		}

		if err != nil {
			if os.IsNotExist(err) {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.Writer.Header().Add("Content-Type", "application/json")

		b = webPathRE.ReplaceAll(b, []byte(dataDirs[i]))

		c.Writer.Header().Add("Content-Type", "text/plain")
		if _, err := io.Copy(c.Writer, bytes.NewBuffer(b)); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})

	// Routes static files. (superfluous if behind properly configured reverse proxy)
	r.NoRoute(func(c *gin.Context) {
		// TODO finish validity check.
		//path := c.Request.URL.Path
		//if !strings.HasPrefix(path, "audio") && !strings.HasPrefix(path, "video") {}

		f, err := os.Open(filepath.Join(".", c.Request.URL.Path))
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

	log.Println(r.Run(":" + strconv.Itoa(*port)))
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
	for _, pd := range playlistDirs {
		err = filepath.Walk(pd, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.Name()[0] == '.' || info.Mode()&0o111 != 0 || info.Name()[len(info.Name())-1] == '~' {
				return nil
			}

			p.Fnames = append(p.Fnames, filepath.Base(info.Name()))
			return nil
		}))
	}
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

	var (
		b  []byte
		i  int // only useful because playlistDirs and dataDirs are aligned.
		pd string
	)
	for i, pd = range playlistDirs {
		if b, err = ioutil.ReadFile(filepath.Join(pd, playlist)); err == nil {
			break
		}
	}
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
			WebPath: webPathRE.ReplaceAllString(fname, dataDirs[i]),
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
