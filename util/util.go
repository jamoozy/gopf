package util

import (
  "log"
  "os"
)

// Utility function that returns true iff `name` exists and is a directory.
func IsDir(name string) bool {
  fileInfo, err := os.Stat(name)
  if err != nil {
    if !os.IsNotExist(err) {
      log.Printf("Warning, got error %s on os.Stat\n", err.Error())
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
      log.Printf("Warning, got error %s on os.Stat\n", err.Error())
    }
    return false
  }
  return !fileInfo.IsDir()
}
