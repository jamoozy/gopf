package dblayer

import (
  "database/sql"
  "errors"
  "flag"
  "fmt"
)

// Internal libraries.
import (
  "github.com/jamoozy/gopf/lg"
  "github.com/jamoozy/gopf/util"
)

// 3rd party libraries.
import (
  // Sqlite3 driver (needed by "sql").
  _ "github.com/mattn/go-sqlite3"
)


// Variable for setting the database path.
var dv = dbVar{path: "gopf.db"}

// The name of the DB file.
type dbVar struct {
  path string
}

// Returns currently set DB path.
func (dv *dbVar) String() string {
  return dv.path
}

// Set the dbVar path.
func (dv *dbVar) Set(path string) error {
  if !util.IsFile(path) {
    return errors.New(fmt.Sprintf("File: %s D.N.E.", path))
  }
  dv.path = path
  return nil
}

// Initialize the database file path with
func init() {
  flag.Var(&dv, "db", "Name of the DB file.")
}



////////////////////////////////////////////////////////////////////////////////
//                                SQL Helpers                                 //
////////////////////////////////////////////////////////////////////////////////

// Type of parsing function for SqlQuery().
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
  for i, v := range orig {
    fnames[i] = v[0]
  }
  return fnames
}

// Function that runs some SQL query.
type SqlRunner func(*sql.DB) error

// Gets a SQL context and passes a *sql.DB to the passed function.
func SqlCtx(fn SqlRunner) error {
  db, err := sql.Open("sqlite3", dv.path)
  if err != nil {
    return err
  }
  defer db.Close()

  return fn(db)
}

// Runs a SQL query.  Parses the rows with the passed function, fn.
func SqlQuery(fn RowParser, stmt string, args ...interface{}) ([][]string, error) {
  var parsedOutput [][]string
  lg.Trc("Running command: %s <-- %s\n", stmt, args)
  return parsedOutput, SqlCtx(func(db *sql.DB) error {
    rows, err := db.Query(stmt, args...)
    if err != nil {
      lg.Ftl(err.Error())
      return err
    }
    defer rows.Close()

    parsedOutput, err = fn(rows)
    return err
  })
}

// Executes a function that doesn't return rows.
func SqlExec(stmt string, args ...interface{}) (error) {
  return SqlCtx(func(db *sql.DB) error {
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
//                       Querying Convenience Functions                       //
////////////////////////////////////////////////////////////////////////////////

// Checks the DB for all the tags.
func QueryTags() ([]string, error) {
  rtn, err := SqlQuery(SingleStringRowParser, "select tags.name from tags")
  if err != nil {
    return nil, err
  }
  return toArray(rtn), nil
}

// Returns list of all files with the given tag.
func QueryFiles(tag string) ([]string, error) {
  rtn, err := SqlQuery(
    SingleStringRowParser,
    "select files.path from files, tags, file_tags" +
    "  where files.id = file_tags.file_id" +
    "    and tags.id = file_tags.tag_id" +
    "    and tags.name = ?", tag)
  if err != nil {
    return nil, err
  }

  return toArray(rtn), nil
}

// Gets all the playlists in the system.
func QueryPlaylists() ([]string, error) {
  playlists, err := SqlQuery(SingleStringRowParser,
                             "select name from playlists")
  if err != nil {
    return nil, err
  }
  return toArray(playlists), nil
}

// Gets all the files in the passed playlist.
func QueryMedia(playlist string) ([]string, error) {
  media, err := SqlQuery(
    SingleStringRowParser,
    `select files.name from playlists, files, playlist_files
       where playlists.id = playlist_files.playlist_id
         and files.id = playlist_files.file_id
         ane playlists.name = ?`, playlist)
  if err != nil {
    return nil, err
  }
  return toArray(media), nil
}



////////////////////////////////////////////////////////////////////////////////
//                       Querying Convenience Functions                       //
////////////////////////////////////////////////////////////////////////////////

// Tags the file with the tag.
func TagFile(tag, file string) error {
  err := SqlExec("insert or ignore into tags('name') values(?)", tag)
  if err != nil {
    return err
  }

  return SqlExec("insert into file_tags(file_id,tag_id)" +
                 "  select files.id, tags.id" +
                 "    from files, tags" +
                 "    where files.path=? and tags.name=?", file, tag)
}
