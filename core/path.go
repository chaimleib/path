package core

import (
  "os"
  "strings"
)

func List(varName string) ([]string, error) {
  result := strings.Split(os.Getenv(varName), ":")
  return result, nil
}

