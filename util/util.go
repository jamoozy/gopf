package util

import (
  "errors"
  "os"
  "strings"
)

import (
  "github.com/jamoozy/gopf/lg"
)

// Utility function that returns true iff `name` exists and is a directory.
func IsDir(name string) bool {
  fileInfo, err := os.Stat(name)
  if err != nil {
    if !os.IsNotExist(err) {
      lg.Wrn("Got error %s on os.Stat\n", err.Error())
    }
    return false
  }
  return fileInfo.IsDir()
}

// Utility function that returns true iff `name` exists and is a file.
func IsFile(name string) bool {
  fileInfo, err := os.Stat(name)
  if err != nil {
    if !os.IsNotExist(err) {
      lg.Wrn("Got error %s on os.Stat\n", err.Error())
    }
    return false
  }
  return !fileInfo.IsDir()
}

// Returns the directory of the specified path.  Can be either a directory or a
// file path.  If the root directory (i.e., "/") is passed, return an error.
func DirOf(path string) (string, error) {
  path = strings.TrimRight(path, "\t ")
  if path == "/" {
    msg := "Attempted to get directory of root directory."
    lg.Trc("%s ... Returning.", msg)
    return "", errors.New(msg)
  }

  lg.Trc("DirOf(%s)", path)
  path = strings.TrimRight(path, "/")
  return path[0:strings.LastIndex(path, "/")], nil
}
