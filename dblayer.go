package dblayer

import (
  "database/sql"
  "errors"
  "log"
)

// 3rd party libraries.
import (
  // Sqlite3 driver (needed by "sql").
  _ "github.com/mattn/go-sqlite3"
)


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



////////////////////////////////////////////////////////////////////////////////
//                                SQL Helpers                                 //
////////////////////////////////////////////////////////////////////////////////

// Function that runs some SQL query.
type SqlRunner func(*sql.DB) error

// Gets a SQL context and passes a *sql.DB to the passed function.
func SqlCtx(fn SqlRunner) error {
  db, err := sql.Open("sqlite3", "sqlite3.db")
  if err != nil {
    return err
  }
  defer db.Close()

  return fn(db)
}

// Runs a SQL query.  Parses the rows with the passed function, fn.
func SqlQuery(fn RowParser, stmt string, args ...interface{}) ([][]string, error) {
  var parsedOutput [][]string
  log.Printf("Running command: %s <-- %s\n", stmt, args)
  return parsedOutput, SqlCtx(func(db *sql.DB) error {
    rows, err := db.Query(stmt, args...)
    if err != nil {
      log.Printf("Error: %s\n", err.Error())
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

// Checks the DB for all the tags.
func QueryTags() ([]string, error) {
  rtn, err := SqlQuery(SingleStringRowParser, "select tags.name from tags")
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


