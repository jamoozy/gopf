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
  "log"
  "net/http"
  "os"
  "regexp"
  "strings"
  "time"
)

// My libraries
import "github.com/jamoozy/gopf/dblayer"


////////////////////////////////////////////////////////////////////////////////
//                                 Endpoints                                  //
////////////////////////////////////////////////////////////////////////////////

// Type that all of my handler functions are.
type handlerFunc func(http.ResponseWriter, *http.Request, ...string) error

// Sets a tag on a file.
func settag(w http.ResponseWriter, r *http.Request, args ...string) error {
  log.Printf("Creating tag for tag:'%s', file:'%s'\n", args)

  err := dblayer.TagFile(args[0], args[1])
  if err != nil {
    return err
  }

  log.Println("Successfully added file tag.")
  return nil
}

// Gets all files tagged with the specified tag.
func gettag(w http.ResponseWriter, r *http.Request, args ...string) error {
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
  log.Println("In /gettags")
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
  logErr := func(err error) {
    log.Println(err)
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
    "^/(gettags|(settag|hastag|gettag)/(.*))/?$")

// Wraps my handlerFunc into an http.HandlerFunc given the set of allowable
// methods and number of additional string arguments to pass to fn.
//
// **Under no circumstances shall `fn` not be passed exactly `numArgs`
// additional arguments.**
func wrapHandler(fn handlerFunc, methods map[string]bool, numArgs int) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    // Ensure this is a valid method for this function.
    if !methods[r.Method] {
      log.Printf("Method not allowed: %s %s", r.Method, r.URL.Path)
      w.WriteHeader(http.StatusMethodNotAllowed)
      return
    }

    // Ensure that the URL matches.
    m := validMethods.FindStringSubmatch(r.URL.Path)
    if m == nil {
      log.Printf("Page '%s' Not found.\n", r.URL.Path)
      http.NotFound(w, r)
      return
    }

    // Split into separate args; make sure there are the right amount.
    args := strings.Split(m[4], "/")
    if len(args) != numArgs {
      http.Error(w, "Wrong #args.", http.StatusBadRequest)
      return
    }

    err := fn(w, r, args...)
    if err != nil {
      log.Println(err)
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
    log.Printf("Satisfied request for index.html")
    http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
    return
  }

  vrb(`rootHandler got request at path "%s"`, r.URL.Path)

  m := servableFiles.FindStringSubmatch(r.URL.Path)
  if m == nil {
    log.Printf(`Unrecognized endpoint: "%s".\n`, r.URL.Path)
    http.NotFound(w, r)
    return
  }
  vrb(`rootHandler serving file: "%s"`, m[1])
  http.ServeFile(w, r, m[1])
}



////////////////////////////////////////////////////////////////////////////////
//                                    Main                                    //
////////////////////////////////////////////////////////////////////////////////

// The GOPF context.
var (
  mediaDir string    // Directory where data is stored.
  port string       // Port to open HTTP(S) server on.
  verbose bool      // Whether to print "verbose" logs.
)

// Print verbosely.
func vrb(fmt string, args ...interface{}) {
  if !verbose {
    return
  }
  log.Printf(fmt, args...)
}

func parseArgs() {
  flag.StringVar(&mediaDir, "data", "data", "Data directory.")
  flag.StringVar(&dblayer.DbName, "db", "gopf.db", "Name of the DB file.")
  flag.StringVar(&port, "port", "8080", "Port to server on.")
  flag.BoolVar(&verbose, "verbose", false, "Switch on verbose mode.")
  flag.Parse()

  // Some minor validation.
  fileInfo, err := os.Stat(mediaDir)
  if err != nil {
    if os.IsNotExist(err) {
      log.Fatalf("%d: does not exist", mediaDir)
    }
    log.Fatalf(err.Error())
  }
  vrb("%s: exists", mediaDir)

  if !fileInfo.IsDir() {
    log.Fatalf("%d: not a directory", mediaDir)
  }
  vrb("%s: is a dir", mediaDir)

  fileInfo, err = os.Stat(dblayer.DbName)
  if err != nil {
    if os.IsNotExist(err) {
      log.Printf("%d: does not exist.  Creating new.", dblayer.DbName)

      // TODO Create new DB.
      log.Fatalln("Not implemented :-(")
    } else {
      vrb("%s: exists")
      log.Fatalln(err.Error())
    }
  }
  vrb("%s: exists", fileInfo.Name())

  if fileInfo.IsDir() {
    log.Fatalln("%s: directory")
  }
  vrb("%s: file", fileInfo.Name())
}

func main() {
  parseArgs()

  get := map[string]bool{"GET": true}
  put := map[string]bool{"PUT": true}

  http.HandleFunc("/settag/", wrapHandler(settag, put, 2))
  http.HandleFunc("/gettag/", wrapHandler(gettag, get, 1))
  http.HandleFunc("/gettags/", wrapHandler(gettags, get, 0))
  http.HandleFunc("/index.html", serveIndex)
  http.HandleFunc("/", rootHandler)

  vrb("Running server on %s", port)
  err := http.ListenAndServe(":" + port, nil)
  if err != nil {
    log.Fatal(err)
  }
}
