package main

import (
  "fmt"
  "log"
  "os"
  "path/core"
  "strings"
)

type Opts struct{
  Commands []*Command
  ExtraArgs []*Arg
  AllArgs []string
}

func NewOpts(osArgs []string) *Opts {
  self := new(Opts)
  self.Commands = make([]*Command, 0)
  self.ExtraArgs = make([]*Arg, 0)
  self.AllArgs = osArgs
  return self
}

func (self *Opts) LastCommand() *Command {
  l := len(self.Commands) - 1
  if l == -1 {
    panic("can't get LastCommand()")
  }
  return self.Commands[l]
}

func (self *Opts) AppendCommand(op string, i int) {
  self.Commands = append(self.Commands, NewCommand(op, i))
}

func (self *Opts) AppendExtraArg(arg string, i int) {
  self.ExtraArgs = append(self.ExtraArgs, NewArg(arg, i))
}

type Arg struct {
  Text string
  Index int
}

func NewArg(text string, i int) *Arg {
  self := new(Arg)
  self.Text = text
  self.Index = i
  return self
}

type Command struct{
  Operation string
  Args []*Arg
  AllArgsIndex int
}

func NewCommand(op string, i int) *Command {
  self := new(Command)
  self.Operation = op
  self.Args = make([]*Arg, 0)
  self.AllArgsIndex = i
  return self
}

// Sorts arguments between Commands (beginning with a "-" and a letter for id)
// and ExtraArgs, not tied to any Command. 0 or 1 ExtraArgs are expected, and
// certain Commands expect no sibling Commands.
func (self *Opts) Parse() error {
  nextLiteral := false
  argsLeft := 0
  for i, arg := range self.AllArgs {
    if nextLiteral {
      nextLiteral = false
      if argsLeft > 0 {
        self.LastCommand().AppendArg(arg, i)
        argsLeft--
      } else {
        self.AppendExtraArg(arg, i)
      }
      continue
    }
    if strings.HasPrefix(arg, "-") {
      if arg == "-" { // bare "-" escapes the next argument
        nextLiteral = true
        continue
      }
      cmd := arg[1:]
      switch cmd[0] {
      case 'l':
        self.AppendCommand(cmd, i)
        cmdArg := cmd[1:]
        if cmdArg != "" {
          self.LastCommand().AppendArg(cmdArg, i)
          argsLeft = 0
        }
      default:
        return fmt.Errorf("unknown flag: %q", arg)
      }
    } else {
      // no '-' prefix
      if argsLeft > 0 {
        self.LastCommand().AppendArg(arg, i)
        argsLeft--
      } else {
        self.AppendExtraArg(arg, i)
      }
    }
  }
  return nil
}

func (self *Command) AppendArg(arg string, i int) {
  self.Args = append(self.Args, NewArg(arg, i))
}

func main() {
  opts := NewOpts(os.Args[1:])
  err := opts.Parse()
  if err != nil {
    log.Fatal(err)
  }
  execute(opts)
}

func execute(opts *Opts) error {
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

