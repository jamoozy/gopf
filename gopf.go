// Copyright (C) 2011-2015 Andrew "Jamoozy" Sabisch
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under the
// terms of the GNU Affero General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// GOPF is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License
// along with GOPF. If not, see http://www.gnu.org/licenses/.


package main

// System libraries
import (
  "bytes"
  "fmt"
  "encoding/json"
  "flag"
  "html/template"
  "net/http"
  "os"
  "path/filepath"
  "regexp"
  "strings"
  "time"

  // Internal libraries.
  "github.com/jamoozy/gopf/dblayer"
  "github.com/jamoozy/gopf/util"
  "github.com/jamoozy/util/lg"
)



////////////////////////////////////////////////////////////////////////////////
//                                 Endpoints                                  //
////////////////////////////////////////////////////////////////////////////////

// Type that all of my handler functions are.
type handlerFunc func(*GopfCall, ...string) error

// These variables together are the GOPF context.
type Gopf struct {
  mediaDir string   // Directory where data is stored.
  port     string   // Port to open HTTP(S) server on.
  wd       string   // Working directory.
  pScan    string   // Where to run dblayer.ScanUpdateDB for playlists.
  tScan    string   // Where to run dblayer.ScanUpdateDB for tags.
  mediaTag string   // Type of HTML tag for media player.
}

type GopfCall struct {
  *Gopf   // Pointer to the GOPF context (read only).

  // The sent response writer.
  w http.ResponseWriter

  // The request.
  r *http.Request
}

func (gopf *Gopf) MakeCall(w http.ResponseWriter, r *http.Request) *GopfCall {
  return &GopfCall{gopf, w, r}
}

// Gets the contents of a playlist.
func (g *GopfCall) playlist(args ...string) error {
  lg.Enter("playlist(%s)\n", args)
  defer lg.Exit("playlist(%s)\n", args)

  files, err := dblayer.GetPlaylist(args[0])
  if err != nil {
    return err
  }

  lg.Vrb("Sending %d files.", len(files))
  rtn := struct{
    Files []string
  }{
    files,
  }

  j, err := json.Marshal(rtn)
  if err != nil {
    return err
  }

  g.ServeContentNow(j)
  return nil
}

func (g *GopfCall) ServeContentNow(content []byte) {
  http.ServeContent(g.w, g.r, "", time.Now().UTC(), bytes.NewReader(content))
}

// Sets a tag on a file.
func (g *GopfCall) settag(args ...string) error {
  lg.Enter("settag(w, r, %s)\n", args)
  defer lg.Exit("settag(w, r, %s)\n", args)

  err := dblayer.TagFile(args[0], args[1])
  if err != nil {
    return err
  }

  lg.Vrb("Successfully added file tag.")
  return nil
}

// Gets all files tagged with the specified tag.
func (g *GopfCall) gettag(args ...string) error {
  lg.Enter("gettag(w, r, %s)\n", args)
  defer lg.Exit("gettag(w, r, %s)\n", args)

  rtn, err := dblayer.QueryFiles(args[0])
  if err != nil {
    return err
  }

  j, err := json.Marshal(rtn)
  if err != nil {
    return err
  }

  http.ServeContent(g.w, g.r, "", time.Now().UTC(), bytes.NewReader(j))
  return nil
}

// Gets a list of all the tags in existence.
func (g *GopfCall) gettags(args ...string) error {
  lg.Enter("gettags(w, r, %s)\n", args)
  defer lg.Exit("gettags(w, r, %s)\n", args)

  rtn, err := dblayer.QueryTags()
  if err != nil {
    return err
  }

  j, err := json.Marshal(rtn)
  if err != nil {
    return err
  }

  http.ServeContent(g.w, g.r, "", time.Now(), bytes.NewReader(j))
  return nil
}

// Simply serves the main page, index.hml
func (g *GopfCall) serveIndex() {
  lg.Enter("serveIndex(w, r)")
  defer lg.Exit("serveIndex(w, r)")

  // Convenience.
  logErr := func(err error) {
    lg.Wrn("Got error: %s", err)
    http.Error(g.w, err.Error(), http.StatusInternalServerError)
  }

  t, err := template.New("index.tmpl").ParseFiles("index.tmpl")
  if err != nil {
    logErr(err)
    return
  }

  type PageTmpl struct {
    Title string
    Playlists []string
    Selected string
    Media []string
    Playing string
    MediaTag template.HTML
  }

  // Build data for the template.
  pt := PageTmpl{
    "GOPF",
    []string{},
    g.r.Form.Get("p"),
    []string{},
    g.r.Form.Get("m"),
    template.HTML(fmt.Sprintf(`<%s id="player" src="" seek="true" controls>Hey, man, get an HTML5-compatible browser, okay?</%s>`, g.mediaTag, g.mediaTag)),
  }

  pt.Playlists, err = dblayer.QueryPlaylists()
  if err != nil {
    logErr(err)
    return
  }
  if pt.Selected != "" {
    pt.Media, err = dblayer.QueryMedia(pt.Selected)
    if err != nil {
      logErr(err)
      return
    }
  }
  err = t.Execute(g.w, pt)
  if err != nil {
    logErr(err)
    return
  }
}



////////////////////////////////////////////////////////////////////////////////
//                         Interface to `http` Module                         //
////////////////////////////////////////////////////////////////////////////////

// Regular expression defining valid URLs.  This variable simplifies the
// redirection process.
var validMethods = regexp.MustCompile(
    "^/(gettags|(settag|hastag|gettag|playlist)/(.*))/?$")

// Wraps my handlerFunc into an http.HandlerFunc given the set of allowable
// methods and number of additional string arguments to pass to fn.
//
// **Under no circumstances shall `fn` not be passed exactly `numArgs`
// additional arguments.**
func (gopf *Gopf) wrapHandler(fn handlerFunc, methods map[string]bool, numArgs int) http.HandlerFunc {
  lg.Enter(`wrapHandler(fn, %v, %d)`, methods, numArgs)
  defer lg.Exit(`wrapHandler(fn, %v, %d)`, methods, numArgs)

  return func(w http.ResponseWriter, r *http.Request) {
    lg.Enter(`wrapHandler_internal(fn, %v, %d).<return>(w, r)`, methods, numArgs)
    defer lg.Exit(`wrapHandler_internal(fn, %v, %d).<return>(w, r)`, methods, numArgs)

    // Ensure this is a valid method for this function.
    if !methods[r.Method] {
      lg.Ifo("Method not allowed: %s %s\n", r.Method, r.URL.Path)
      w.WriteHeader(http.StatusMethodNotAllowed)
      return
    }

    // Ensure that the URL matches.
    m := validMethods.FindStringSubmatch(r.URL.Path)
    if m == nil {
      lg.Ifo("Page '%s' Not found.\n", r.URL.Path)
      http.NotFound(w, r)
      return
    }

    // Split into separate args; make sure there are the right amount.
    args := strings.Split(m[3], "/")
    if len(args) != numArgs {
      msg := "Wrong #args."
      lg.Ifo(msg)
      http.Error(w, msg, http.StatusBadRequest)
      return
    }

    // Call the function and wrap any errors.
    err := fn(gopf.MakeCall(w, r), args...)
    if err != nil {
      lg.Ifo(err.Error())
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }
}

// Matches HTML, JavaScript, and CSS files for the default handler.
var staticFiles = regexp.MustCompile(`/([a-zA-Z0-9_.-]+\.(html|js|css))`)

// Matches all media files.
var mediaFiles = regexp.MustCompile(`/(.*\.(mp[34]|ogg|ogv))`)

// Handles all default requests.
func (g *GopfCall) rootHandler() {
  lg.Enter("rootHandler(w, r)")
  defer lg.Exit("rootHandler(w, r)")

  if g.r.URL.Path == "" || g.r.URL.Path == "/" {
    lg.Ifo("Satisfied request for index.html")
    http.Redirect(g.w, g.r, "/index.html", http.StatusMovedPermanently)
    return
  }

  lg.Trc(`rootHandler got request at path "%s"`, g.r.URL.Path)

  m := staticFiles.FindStringSubmatch(g.r.URL.Path)
  if m != nil {
    g.serveFile("static" + g.r.URL.Path)
    return
  }

  m = mediaFiles.FindStringSubmatch(g.r.URL.Path)
  if m != nil {
    g.serveFile(g.mediaDir + g.r.URL.Path)
    return
  }
  lg.Ifo(`Unrecognized endpoint: "%s".`, g.r.URL.Path)
  http.NotFound(g.w, g.r)
}

// Serves a file, or a "404 Not Found".
func (g *GopfCall) serveFile(path string) {
  lg.Trc(`serveFile(w, r, "%s")`, path)
  if util.IsFile(path) {
    http.ServeFile(g.w, g.r, path)
  } else {
    lg.Dbg(`404 Not Found`, path)
    http.NotFound(g.w, g.r)
  }
}



////////////////////////////////////////////////////////////////////////////////
//                                    Main                                    //
////////////////////////////////////////////////////////////////////////////////

// Exit error codes.
const (
  DbDne        = -iota
  MediaDirDne  = -iota
  BadMediaType = -iota
)

// Determines what the media tag should be (audio or video) based on the file
// types.
func (g *Gopf) determineMediaTag() {
  // Determine what kind of tag is most appropriate -- video or audio.
  var (
    audio = 0
    video = 0

    // TODO think of more file types or find a library that has some kind of
    //      recognition capabilities ...
    audioRegexp = regexp.MustCompile(`.*\.(mp3|wav|ogg)$`)
    videoRegexp = regexp.MustCompile(`.*\.(mp4|ogv|wmv)$`)
  )

  // The traversal function -- just updates `audio` and `video` to reflect the
  // number of files we've seen.
  traverse := func(path string, info os.FileInfo, err error) error {
    name := info.Name()
    // Directory -- not relevant for determining file type.
    if info.IsDir() {
      // This is a symlink.  We can't handle these because `os.Walk()` doesn't
      // follow symlinks (see https://golang.org/pkg/path/filepath/#Walk).
      // Report that there's an issue.
      if (info.Mode() & os.ModeSymlink) != 0 {
        lg.Wrn(`Can't traverse symlink: "%s"`, name)
      }
      return nil
    }

    // Check which regular expression matches, but favor audio over video.
    if audioRegexp.FindStringSubmatch(name) != nil {
      audio += 1
    } else if videoRegexp.FindStringSubmatch(name) != nil {
      video += 1
    } else {
      lg.Wrn(`Unrecognized file type: "%s"`, name)
    }

    return nil
  }

  // Do the traversal, find the most common type of file, and set the media tag
  // based on majority count.  (note that the default value is "audio")
  err := filepath.Walk(g.mediaDir, traverse)
  if err != nil {
    lg.Wrn(err.Error())
  }
  lg.Dbg("Got %d audio vs. %d video.", audio, video)
  if video > audio {
    g.mediaTag = "video"
  }
  lg.Dbg(`Decided to use media tag: <%s>`, g.mediaTag)
}

func main() {
  gopf := &Gopf{}

  // Parse command-line arguments.
  flag.StringVar(&gopf.mediaDir, "media", "media", "Data directory.")
  flag.StringVar(&gopf.port, "port", "8080", "Port to server on.")
  flag.StringVar(
    &gopf.pScan, "p-scan", "",
    `Scans the directory and populates the DB with playlists based on its
     structure.`)
  flag.StringVar(
    &gopf.tScan, "t-scan", "",
    `Scans the directory and populates the DB with tags based on its
     structure.`)
  flag.StringVar(&gopf.mediaTag, "type", "", "Set media type: audio or video.")
  flag.Parse()

  if err := dblayer.VerifyDb(); err != nil {
    lg.Ftl(err.Error())
    lg.Ftl("  Set db with -db [name]")
    os.Exit(DbDne)
  }

  // Some minor validation.
  if !util.IsDir(gopf.mediaDir) {
    lg.Err(`Media directory: "%s" is not a directory`, gopf.mediaDir)
    lg.Err(`  To set it, run with -media=[file]`)
    os.Exit(MediaDirDne)
  }

  // TODO implement "-type" flag.
  if gopf.mediaTag == "" {
    gopf.determineMediaTag()
  } else if gopf.mediaTag != "audio" && gopf.mediaTag != "video" {
    lg.Ftl(`Invalid media tag type: "%s"`, gopf.mediaTag)
    lg.Ftl(`  expected "audio" or "video".`)
    os.Exit(BadMediaType)
  }

  // Update the DB if it was requested to do so.
  if gopf.pScan != "" {
    dblayer.ScanUpdate(gopf.pScan, `playlists`)
  }
  if gopf.tScan != "" {
    dblayer.ScanUpdate(gopf.tScan, `tags`)
  }
  if gopf.pScan != "" || gopf.tScan != "" {
    // Exit if we did a pScan and/or tScan
    os.Exit(0)
  }

  // Establish working directory.
  var err error
  gopf.wd, err = os.Getwd()
  if err != nil {
    lg.Wrn(`Can't determine working directory.`)
    gopf.wd = `.`
  }
  lg.Vrb("Running server on %s at %s", gopf.port, gopf.wd)

  // Convenience sets.
  get := map[string]bool{"GET": true}
  put := map[string]bool{"PUT": true}

  // All the endpoints.
  http.HandleFunc("/playlist/", gopf.wrapHandler(func(g *GopfCall, args ...string) error {
    return g.playlist(args...)
  }, get, 1))
  http.HandleFunc("/settag/", gopf.wrapHandler(func(g *GopfCall, args ...string) error {
    return g.settag(args...)
  }, put, 2))
  http.HandleFunc("/gettag/", gopf.wrapHandler(func(g *GopfCall, args ...string) error {
    return g.gettag(args...)
  }, get, 1))
  http.HandleFunc("/gettags/", gopf.wrapHandler(func(g *GopfCall, args ...string) error {
    return g.gettags(args...)
  }, get, 0))
  http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
    gopf.MakeCall(w, r).serveIndex()
  })
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    gopf.MakeCall(w, r).rootHandler()
  })

  // Run and report any shutdown errors.
  err = http.ListenAndServe(":" + gopf.port, nil)
  if err != nil {
    lg.Ftl(err.Error())
  }
}
