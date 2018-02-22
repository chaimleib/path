package opts

import (
  "fmt"
  "strings"
)

type Opts struct{
  Menu []CommandSpec
  Commands []*Command
  ExtraArgs []*Arg
  AllArgs []string
}

type CommandSpec interface {
  Letters() []rune
  Words()   []string
  Help()    string
  ExpectedArgs() int
  ArgsRequired() bool
}

func New(osArgs []string, menu []CommandSpec) *Opts {
  self := new(Opts)
  self.Menu = menu
  self.Commands = make([]*Command, 0)
  self.ExtraArgs = make([]*Arg, 0)
  self.AllArgs = osArgs
  return self
}

func (self *Opts) MenuLetter(l rune) CommandSpec {
  for _, spec := range self.Menu {
    for _, specLetter := range spec.Letters() {
      if l == specLetter {
        return spec
      }
    }
  }
  return nil
}

func (self *Opts) MenuWord(w string) CommandSpec {
  for _, spec := range self.Menu {
    for _, specWord := range spec.Words() {
      if specWord == w {
        return spec
      }
    }
  }
  return nil
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

