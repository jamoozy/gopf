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
  "database/sql"
  "encoding/json"
  "errors"
  "fmt"
  "html/template"
  "log"
  "net/http"
  "regexp"
//  "strings"
  "time"
)

// My libraries
import "libs/dblayer"


////////////////////////////////////////////////////////////////////////////////
//                                 Endpoints                                  //
////////////////////////////////////////////////////////////////////////////////

// Type that all of my handler functions are.
type handlerFunc func(http.ResponseWriter, *http.Request, string, string) error

// Sets a tag on a file.
func settag(w http.ResponseWriter, r *http.Request, tag string, file string) error {
  log.Printf("Creating tag for tag:'%s', file:'%s'\n", tag, file)

  err := dblayer.SqlExec("insert or ignore into tags('name') values(?)", tag)
  if err != nil {
    return err
  }

  err = dblayer.SqlExec("insert into file_tags(file_id,tag_id)" +
                        "  select files.id, tags.id" +
                        "    from files, tags" +
                        "    where files.path=? and tags.name=?", file, tag)
  if err != nil {
    return err
  }

  log.Println("Successfully added file tag.")
  return nil
}

// Gets all files in the tag.  "file" is ignored; it's only there so this
// conforms to the handlerFunc type.
func gettag(w http.ResponseWriter, r *http.Request, tag string, file string) error {
  rtn, err := dblayer.SqlQuery(
    SingleStringRowParser,
    "select files.path from files, tags, file_tags" +
    "  where files.id = file_tags.file_id" +
    "    and tags.id = file_tags.tag_id" +
    "    and tags.name = ?", tag)
  if err != nil {
    return err
  }

  j, err := json.Marshal(toArray(rtn))
  if err != nil {
    return err
  }

  http.ServeContent(w, r, "", time.Now(), bytes.NewReader(j))
  return nil
}

// Gets a list of all the tags in existence.
func gettags(w http.ResponseWriter, r *http.Request, tag string, file string) error {
  log.Println("In /gettags")
  rtn, err := dblayer.QueryTags()
  if err != nil {
    return err
  }

  fmt.Printf("There are %d return values.\n", len(rtn))

  j, err := json.Marshal(toArray(rtn))
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
    "^/(gettags|(settag|hastag|gettag)/([^/]+)/([^/]+))/?$")

// Wraps my handlerFunc into an http.HandlerFunc given the set of allowable
// methods.
func wrapHandler(fn handlerFunc, methods map[string]bool) http.HandlerFunc {
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

    err := fn(w, r, m[3], m[4])
    if err != nil {
      log.Println(err)
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }
}

// Matches HTML, JavaScript, and CSS files for the default handler.
var servableFiles = regexp.MustCompile(
  "/(.*\\.(html|js|css)|data/.*\\.(mp[34]|ogg|ogv))")

// Handles all default requests.
func rootHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path == "/" {
    log.Printf("Satisfied request for index.html")
    http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
    return
  }

  m := servableFiles.FindStringSubmatch(r.URL.Path)
  if m == nil {
    log.Printf(`Unrecognized endpoint: "%s".\n`, r.URL.Path)
    http.NotFound(w, r)
    return
  }
  log.Println("Got request for file: " + m[1])
  http.ServeFile(w, r, m[1])
}



////////////////////////////////////////////////////////////////////////////////
//                                    Main                                    //
////////////////////////////////////////////////////////////////////////////////

func main() {
  get := map[string]bool{"GET": true}
  put := map[string]bool{"PUT": true}

  http.HandleFunc("/settag/", wrapHandler(settag, put))
  http.HandleFunc("/gettag/", wrapHandler(gettag, get))
  http.HandleFunc("/gettags/", wrapHandler(gettags, get))
  http.HandleFunc("/index.html", serveIndex)
  http.HandleFunc("/", rootHandler)

  port := ":8079"
  fmt.Println("Running server on " + port)
  http.ListenAndServe(port, nil)
}
