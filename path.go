package main

import (
  "fmt"
  "log"
  "os"
  "path/core"
  "strings"
)

type Opts struct{
  VarName string
  Commands []Command
}

type Command struct{
  Operation string
  Arg string
}

func main() {
  opts, err := parseOpts()
  if err != nil {
    log.Fatal(err)
  }
  execute(opts)
}

func parseOpts() (*Opts, error) {
  opts := new(Opts)
  opts.Commands = make([]Command, 0)
  nonFlags := make([]string, 0)
  nextLiteral := false
  setLastArg := false
  for _, arg := range os.Args[1:] {
    if nextLiteral {
      if setLastArg {
        lastCmd := opts.Commands[len(opts.Commands)-1]
        lastCmd.Arg = arg
        setLastArg = false
      } else {
        // parse as nonFlag
        nonFlags = append(nonFlags, arg)
      }
      nextLiteral = false
      continue
    }
    if strings.HasPrefix(arg, "-") {
      if arg == "-" {
        nextLiteral = true
        continue
      }
      cmd := arg[1:]
      switch cmd[0] {
      case 'l':
        opts.Commands = append(opts.Commands, Command{
          Operation: cmd,
          Arg: cmd[1:],
        })
        if cmd[1:] == "" {
          setLastArg = true
        }
      default:
        return nil, fmt.Errorf("unknown flag: %q", arg)
      }
    }
    // no '-' prefix
    if setLastArg {
      lastCmd := opts.Commands[len(opts.Commands)-1]
      lastCmd.Arg = arg
      setLastArg = false
    } else {
      nonFlags = append(nonFlags, arg)
    }
  }
  switch len(nonFlags) {
  case 0:
    opts.VarName = "PATH"
  case 1:
    opts.VarName = nonFlags[0]
  default:
    return nil, fmt.Errorf(
      "expected only one variable name, got %d", len(nonFlags))
  }
  return opts, nil
}

func execute(opts *Opts) error {
  for _, cmd := range opts.Commands {
    switch cmd.Operation {
    case "l":
      paths, err := core.List(opts.VarName)
      if err != nil {
        log.Fatal(err)
      }
      pathLines := strings.Join(paths, "\n")
      fmt.Println(pathLines)
    default:
      return fmt.Errorf("unknown command: %q", cmd)
    }
  }
  return nil
}

