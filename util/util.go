package util

import (
  "os"
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
