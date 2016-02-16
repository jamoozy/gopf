package main

import (
  "database/sql"
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"

  // 3rd party libraries.
  _ "github.com/mattn/go-sqlite3"

  // Internal libraries.
  "github.com/jamoozy/gopf/util"
  "github.com/jamoozy/util/lg"
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
  lg.Enter(`Set("%s")`, path)
  defer lg.Exit(`Set("%s")`, path)

  if !util.IsFile(path) {
    return fmt.Errorf("File: %s D.N.E.", path)
  }
  dv.path = path
  return nil
}

// Initialize the database file path with
func init() {
  flag.Var(&dv, "db", "Name of the DB file.")
}

// Verifies that the database exists, is a file, and is actually a sqlite3 DB.
func VerifyDb() error {
  if !util.IsExists(dv.path) {
    return fmt.Errorf(`Database does not exist: "%s"`, dv.path)
  } else if !util.IsFile(dv.path) {
    return fmt.Errorf(`Database is not a file: "%s"`, dv.path)
  }
  return nil
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
      return fmt.Errorf("No rows affected.")
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
  lg.Trc("QueryMedia(%s)", playlist)

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

// Gets the list of file names in the playlist.
func GetPlaylist(playlist string) ([]string, error) {
  lg.Trc("GetPlaylist(%s)", playlist)

  files, err := SqlQuery(
    SingleStringRowParser,
    `select files.path from files, playlists, playlist_files
       where playlist_files.playlist_id = playlists.id
         and playlist_files.file_id = files.id
         and playlists.name = ?`, playlist)
  if err != nil {
    return nil, err
  }
  return toArray(files), nil
}

// Tags the file with the tag.
func TagFile(tag, file string) error {
  lg.Trc("Tagging file %s with %s", file, tag)

  err := SqlExec(`insert or ignore into tags(name) values(?)`, tag)
  if err != nil {
    lg.Trc("Error: %s", err.Error())
    return err
  }

  return SqlExec(`insert into file_tags(file_id,tag_id)
                    select files.id, tags.id
                      from files, tags
                      where files.path=? and tags.name=?`, file, tag)
}

func TagFiles(tag string, files ...string) error {
  lg.Trc(`TagFiles(%s, %d files)`, tag, len(files))
  for _, file := range files {
    TagFile(tag, file)
  }
  return nil
}

func AddToPlaylist(playlist, file string) error {
  lg.Trc("AddToPlaylist(%s, %s)", playlist, file)
  return SqlExec(`insert into playlist_files(file_id,playlist_id)
                    select files.id, playlists.id
                      from files, playlists
                      where files.path=? and playlists.name=?`, file, playlist)
}

func CreatePlaylist(playlist string, fPaths ...string) (err error) {
  lg.Trc("Creating playlist %s with %d entries.", playlist, len(fPaths))

  err = SqlExec("insert or fail into playlists(name) values(?)", playlist)
  if err != nil {
    return err
  }

  for _, fPath := range fPaths {
    err = SqlExec("insert or ignore into files(path) values(?)", fPath)
    if err != nil {
      lg.Trc("Error inserting %s: %s", fPath, err.Error())
      return err
    }
    err = AddToPlaylist(playlist, fPath)
    if err != nil {
      lg.Trc("Error adding %s to %s: %s", fPath, playlist, err.Error())
      return err
    }
  }

  lg.Trc("Done creating playlist %s", playlist)
  return nil
}



////////////////////////////////////////////////////////////////////////////////
//                            Management Functions                            //
////////////////////////////////////////////////////////////////////////////////

// Context for seeding a DB.
type SeedCtx struct {
  dbPath string       // Path to Database file.
  mediaDir string     // Path to media.  If empty, does not populate the DB,
                      // merely creates it and sets up the schema.
  schemaPath string   // Path to schema file.
  overwrite bool      // Whether to overwrite the DB if it already exists.  If
                      // the DB exists and this is false (default), then Run()
                      // returns an error.
}

// Run a seed process.  **Must not be asynchronous.**
func (sc *SeedCtx) Run() error {
  lg.Trc("SeedCtx.Run(): %s", sc)

  // Default DB path.
  if sc.dbPath == "" {
    sc.dbPath = "gopf.db"
    lg.Dbg("Setting default dbPath: %s", sc.dbPath)
  }

  // Default media dir.
  if sc.mediaDir == "" {
    sc.mediaDir = "media"
    lg.Dbg("Setting default mediaDir: %s", sc.mediaDir)
  }

  // Default schema path.
  if sc.schemaPath == "" {
    sc.mediaDir = "schema.sql"
    lg.Dbg("Setting default schemaPath: %s", sc.schemaPath)
  }

  // Verify that the database path is valid.  This consists of verifying either:
  //  1) It does not exist, should be created, and its directory exists.
  //  2) It exists and should be overwritten.
  if dir := filepath.Dir(sc.dbPath) ; !util.IsDir(sc.dbPath) {
    msg := fmt.Sprintf("No such directory: %s", dir)
    lg.Ifo(msg)
    return fmt.Errorf(msg)
  } else if !sc.overwrite && util.IsFile(sc.dbPath) {
    msg := "File exists, and set not to overwrite."
    lg.Ifo(msg)
    return fmt.Errorf(msg)
  }

  if sc.dbPath != "" {
    // Temporarily swap out current path so that we can use all the nice
    // convenience functions, e.g., SqlCtx.
    var oldName string
    oldName, dv.path = dv.path, sc.dbPath
    defer func() { dv.path = oldName }()
  }

  return nil
}

// This function runs a search on the media directory and populates the DB with
// any missing entries.
func ScanUpdate(mediaDir, table string) error {
  lg.Trc("ScanUpdate(%s, %s)", mediaDir, table)
  _, err := handleDir(mediaDir, table)
  if err != nil {
    lg.Trc("exiting ScanUpdate(%s, %s) err: %s", mediaDir, table, err.Error())
  } else {
    lg.Trc("exiting ScanUpdate(%s, %s) complete.", mediaDir, table)
  }
  return err
}

// Handles the directory.
func handleDir(dirPath, table string) ([]os.FileInfo, error) {
  lg.Trc("Checking out dir: %s", dirPath)

  // Get entries.
  fis, err := ioutil.ReadDir(dirPath)
  if err != nil {
    lg.Trc("Error from ioutil.ReadDir: %s", err.Error())
    return nil, err
  }

  // Split all entries into files and directories.
  dFIs, fFIs := make([]os.FileInfo, 0, 100), make([]os.FileInfo, 0, 100)
  for _, fi := range fis {
    if fi.IsDir() {
      dFIs = append(dFIs, fi)
    } else {
      fFIs = append(fFIs, fi)
    }
  }
  lg.Dbg("%d dirs, %d files", len(dFIs), len(fFIs))

  // Dive into directories, adding their files to the files in this.
  for _, dPath := range dFIs {
    fis, err := handleDir(dPath.Name(), table)
    if err != nil {
      return nil, err
    }
    fFIs = append(fFIs, fis...)
  }
  lg.Dbg("Now have %d files.", len(fFIs))

  // Add all the files into a playlist for this directory.
  fPaths := make([]string, len(fFIs))
  for i, fFI := range fFIs {
    fPaths[i] = fFI.Name()
  }
  switch table {
  case "playlists":
    return fFIs, CreatePlaylist(filepath.Base(dirPath), fPaths...)
  case "tags":
    return fFIs, TagFiles(filepath.Base(dirPath), fPaths...)
  }
  return nil, fmt.Errorf(`No such table: %s`, table)
}
