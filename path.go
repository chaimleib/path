package main

import (
  "fmt"
  "log"
  "os"
  "path/core"
  "path/opts"
  "strings"
)

type Spec struct {
  letters string
  words   []string
  help string
  minArgs int
  maxArgs int
}

func (self Spec) Letters() []rune {
  return []rune(self.letters)
}

func (self Spec) Words() []string {
  return self.words
}

func (self Spec) Help() string {
  return self.help
}

func (self Spec) ExpectedArgs(args []*opts.Arg) int {
  return self.maxArgs - len(args)
}

func (self Spec) ArgsRequired(args []*opts.Arg) bool {
  return len(args) < self.minArgs
}

func main() {
  opts := opts.New(os.Args[1:], []opts.CommandSpec{
    Spec{"l", []string{"list"}, "list members line-by-line", 0, 0},
  }) // TODO: use []CommandSpec
  err := opts.Parse()
  if err != nil {
    log.Fatal(err)
  }
  if len(opts.Commands) == 0 {
    log.Fatal("no commands found")
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
    switch cmd.Text {
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
      return fmt.Errorf("unknown command: %q", opts.AllArgs[cmd.Index])
    }
  }
  return nil
}

