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
  "encoding/json"
  "flag"
  "fmt"
  "html/template"
  "net/http"
  "os"
  "regexp"
  "strings"
  "time"
)

// Internal libraries.
import (
  "github.com/jamoozy/gopf/dblayer"
  "github.com/jamoozy/gopf/util"
  "github.com/jamoozy/gopf/lg"
)


////////////////////////////////////////////////////////////////////////////////
//                                 Endpoints                                  //
////////////////////////////////////////////////////////////////////////////////

// Type that all of my handler functions are.
type handlerFunc func(http.ResponseWriter, *http.Request, ...string) error

// Gets the contents of a playlist.
func playlist(w http.ResponseWriter, r *http.Request, args ...string) error {
  lg.Trc("playlist(w, r, %s)\n", args)

  files, err := dblayer.GetPlaylist(args[0])
  if err != nil {
    return err
  }

  lg.Vrb("Sending %d files.", len(files))
  rtn := struct {
    Files []string
  }{
    Files: files,
  }

  j, err := json.Marshal(rtn)
  if err != nil {
    return err
  }

  http.ServeContent(w, r, "", time.Now(), bytes.NewReader(j))
  return nil
}

// Sets a tag on a file.
func settag(w http.ResponseWriter, r *http.Request, args ...string) error {
  lg.Trc("settag(w, r, %s)\n", args)

  err := dblayer.TagFile(args[0], args[1])
  if err != nil {
    return err
  }

  lg.Vrb("Successfully added file tag.")
  return nil
}

// Gets all files tagged with the specified tag.
func gettag(w http.ResponseWriter, r *http.Request, args ...string) error {
  lg.Trc("gettag(w, r, %s)\n", args)

  rtn, err := dblayer.QueryFiles(args[0])
  if err != nil {
    return err
  }

  j, err := json.Marshal(rtn)
  if err != nil {
    return err
  }

  http.ServeContent(w, r, "", time.Now(), bytes.NewReader(j))
  return nil
}

// Gets a list of all the tags in existence.
func gettags(w http.ResponseWriter, r *http.Request, args ...string) error {
  lg.Trc("gettags(w, r, %s)\n", args)
  rtn, err := dblayer.QueryTags()
  if err != nil {
    return err
  }

  j, err := json.Marshal(rtn)
  if err != nil {
    return err
  }

  http.ServeContent(w, r, "", time.Now(), bytes.NewReader(j))
  return nil
}

// Simply serves the main page, index.hml
func serveIndex(w http.ResponseWriter, r *http.Request) {
  lg.Trc("serveIndex(w, r)")

  // Convenience.
  logErr := func(err error) {
    lg.Wrn("Got error: %s", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  t, err := template.New("index.tmpl").ParseFiles("index.tmpl")
  if err != nil {
    logErr(err)
    return
  }

  // Build data for the template.
  data := struct {
    Title string
    Playlists []string
    Selected string
    Media []string
    Playing string
  }{
    "GOPF",
    []string{},
    r.Form.Get("p"),
    []string{},
    r.Form.Get("m"),
  }

  data.Playlists, err = dblayer.QueryPlaylists()
  if err != nil {
    logErr(err)
    return
  }
  if data.Selected != "" {
    data.Media, err = dblayer.QueryMedia(data.Selected)
    if err != nil {
      logErr(err)
      return
    }
  }
  err = t.Execute(w, data)
  if err != nil {
    logErr(err)
    return
  }
}



////////////////////////////////////////////////////////////////////////////////
//                          Interface to http Module                          //
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
func wrapHandler(fn handlerFunc, methods map[string]bool, numArgs int) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
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
    err := fn(w, r, args...)
    if err != nil {
      lg.Ifo(err.Error())
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }
}

// Matches HTML, JavaScript, and CSS files for the default handler.
var servableFiles = regexp.MustCompile(
    fmt.Sprintf("/([a-zA-Z0-9_.-]+\\.(html|js|css)|media/.*\\.(mp[34]|ogg|ogv))"))

// Handles all default requests.
func rootHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path == "" || r.URL.Path == "/" {
    lg.Ifo("Satisfied request for index.html")
    http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
    return
  }

  lg.Trc(`rootHandler got request at path "%s"`, r.URL.Path)

  m := servableFiles.FindStringSubmatch(r.URL.Path)
  if m == nil {
    lg.Ifo(`Unrecognized endpoint: "%s".\n`, r.URL.Path)
    http.NotFound(w, r)
    return
  }

  if fname := "static" + m[0] ; util.IsFile(fname) {
    http.ServeFile(w, r, fname)
  } else {
    lg.Dbg("%s not found.", fname)
    http.NotFound(w, r)
  }
}



////////////////////////////////////////////////////////////////////////////////
//                                    Main                                    //
////////////////////////////////////////////////////////////////////////////////

// These variables together are the GOPF context.
var (
  mediaDir string         // Directory where data is stored.
  port string             // Port to open HTTP(S) server on.
  wd string               // Working directory.
  shouldScanUpdateDb bool // Whether to run dblayer.ScanUpdateDB
)

// Set default, parse, and validate args.
func parseArgs() {
  flag.StringVar(&mediaDir, "media", "media", "Data directory.")
  flag.StringVar(&port, "port", "8080", "Port to server on.")
  flag.BoolVar(&shouldScanUpdateDb, "scan", false,
               "Scans the media directory and populates the DB with playlists" +
               " based on directory structure.")
  flag.Parse()

  // Some minor validation.
  util.IsFile(mediaDir)
}

func main() {
  parseArgs()

  // Update the DB if it was requested to do so.
  if shouldScanUpdateDb {
    dblayer.ScanUpdateDB(mediaDir)
    return
  }

  // Establish working directory.
  var err error
  wd, err = os.Getwd()
  if err != nil {
    lg.Wrn("can't determine working directory.")
    wd = "."
  }
  lg.Vrb("Running server on %s at %s", port, wd)

  // Convenience sets.
  get := map[string]bool{"GET": true}
  put := map[string]bool{"PUT": true}

  // All the endpoints.
  http.HandleFunc("/playlist/", wrapHandler(playlist, get, 1))
  http.HandleFunc("/settag/", wrapHandler(settag, put, 2))
  http.HandleFunc("/gettag/", wrapHandler(gettag, get, 1))
  http.HandleFunc("/gettags/", wrapHandler(gettags, get, 0))
  http.HandleFunc("/index.html", serveIndex)
  http.HandleFunc("/", rootHandler)

  // Run and report any shutdown errors.
  err = http.ListenAndServe(":" + port, nil)
  if err != nil {
    lg.Ftl(err.Error())
  }
}
