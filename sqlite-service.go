package main

// System libraries
import (
  "bytes"
  "database/sql"
  "encoding/json"
  "errors"
  "fmt"
  "log"
  "net/http"
  "regexp"
  "strings"
  "time"
)

// 3rd party libraries.
import (
  _ "github.com/mattn/go-sqlite3"
)

// Type of parsing function for sqlQuery().
type RowParser func(*sql.Rows) ([][]string, error)

// Creates a simple RowParser that returns an error if rows is empty.
func SingleStringRowParser(rows *sql.Rows) ([][]string, error) {
  rtn := make([][]string, 0, 100)  // 100 is a guess.
  for rows.Next() {
    var fname string
    if err := rows.Scan(&fname); err != nil {
      return nil, err
    }
    rtn = append(rtn, []string{fname})
  }
  return rtn, nil
}

// Converts a [][]string with one string per sub-array into a simple []string
// with the single string contents of each sub-[]string from the original
// [][]string as an entry.
func toArray(orig [][]string) []string {
  fnames := make([]string, len(orig))
  for i := 0; i < len(orig); i++ {
    fnames[i] = orig[i][0]
  }
  return fnames
}



////////////////////////////////////////////////////////////////////////////////
//                                SQL Helpers                                 //
////////////////////////////////////////////////////////////////////////////////

// Gets a SQL context and passes a *sql.DB to the passed function.
func sqlCtx(fn func(*sql.DB) error) error {
  db, err := sql.Open("sqlite3", "sqlite3.db")
  if err != nil {
    return err
  }
  defer db.Close()

  return fn(db)
}

// Runs a SQL query.  Parses the rows with the passed function, fn.
func sqlQuery(fn RowParser, stmt string, args ...interface{}) ([][]string, error) {
  var parsedOutput [][]string
  log.Printf("Running command: %s <-- %s\n", stmt, args)
  return parsedOutput, sqlCtx(func(db *sql.DB) error {
    rows, err := db.Query(stmt, args...)
    if err != nil {
      return err
    }
    defer rows.Close()

    parsedOutput, err = fn(rows)
    return err
  })
}

// Executes a function that doesn't return rows.
func sqlExec(stmt string, args ...interface{}) (error) {
  return sqlCtx(func(db *sql.DB) error {
    result, err := db.Exec(stmt, args...)
    if err != nil {
      return err
    }

    num, err := result.RowsAffected()
    if err != nil {
      return err
    } else if num > 0 {
      return nil
    } else {
      return errors.New("No rows affected.")
    }
  })
}



////////////////////////////////////////////////////////////////////////////////
//                                 Endpoints                                  //
////////////////////////////////////////////////////////////////////////////////

// Type that all of my handler functions are.
type handlerFunc func(http.ResponseWriter, *http.Request, string, string)

// Sets a tag on a file.
func settag(w http.ResponseWriter, r *http.Request, tag string, file string) {
  log.Printf("Creating tag for tag:'%s', file:'%s'\n", tag, file)

  err := sqlExec("insert or ignore into tags('name') values(?)", tag)
  if err != nil {
    log.Println(err)
    http.ServeContent(w, r, "", time.Now(), strings.NewReader(err.Error()))
    return
  }

  err = sqlExec("insert into file_tags(file_id,tag_id)" +
                "  select files.id, tags.id" +
                "    from files, tags" +
                "    where files.path=? and tags.name=?", file, tag)
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusBadRequest)
    http.ServeContent(w, r, "", time.Now(), strings.NewReader("No such file."))
    return
  }

  log.Println("Successfully added file tag.")
}

// Gets all files in the tag.  "file" is ignored; it's only there so this
// conforms to the handlerFunc type.
func gettag(w http.ResponseWriter, r *http.Request, tag string, file string) {
  rtn, err := sqlQuery(
    SingleStringRowParser,
    "select files.path from files, tags, file_tags" +
    "  where files.id = file_tags.file_id" +
    "    and tags.id = file_tags.tag_id" +
    "    and tags.name = ?", tag)
  if err != nil {
    log.Println(err)
    return
  }

  j, err := json.Marshal(toArray(rtn))
  if err != nil {
    log.Println(err)
    return
  }

  http.ServeContent(w, r, "", time.Now(), bytes.NewReader(j))
}

func gettags(w http.ResponseWriter, r *http.Request, tag string, file string) {
  log.Println("In /gettags")
  rtn, err:= sqlQuery(SingleStringRowParser, "select tags.name from tags")
  if err != nil {
    log.Println(err)
    return
  }

  fmt.Printf("There are %d return values.\n", len(rtn))

  j, err := json.Marshal(toArray(rtn))
  if err != nil {
    log.Println(err)
    return
  }

  http.ServeContent(w, r, "", time.Now(), bytes.NewReader(j))
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

    fn(w, r, m[3], m[4])
  }
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

  fmt.Println("Running server on port 8079")
  http.ListenAndServe(":8079", nil)
}
