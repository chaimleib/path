package main

import (
  "fmt"
  "log"
  "os"
  "path/core"
  "path/opts"
  "strings"
)

func main() {
  opts := opts.New(os.Args[1:], nil) // TODO: use []CommandSpec
  err := opts.Parse()
  if err != nil {
    log.Fatal(err)
  }
  execute(opts)
}

func execute(opts *opts.Opts) error {
  pathVar := "PATH" // default
  if len(opts.ExtraArgs) == 1 {
    pathVar = opts.ExtraArgs[0].Text
  } else if len(opts.ExtraArgs) != 0 {
    return fmt.Errorf(
      "only one path variable is allowed, found more than one: %q",
      opts.ExtraArgs)
  }
  for _, cmd := range opts.Commands {
    switch cmd.Operation {
    case "l":
      if len(opts.Commands) > 1 {
        return fmt.Errorf("-l can't be with other commands")
      }
      paths, err := core.List(pathVar)
      if err != nil {
        return err
      }
      pathLines := strings.Join(paths, "\n")
      fmt.Println(pathLines)
    default:
      return fmt.Errorf("unknown command: %q", opts.AllArgs[cmd.AllArgsIndex])
    }
  }
  return nil
}

