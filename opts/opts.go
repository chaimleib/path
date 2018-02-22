package opts

import (
  "fmt"
  "strings"
  "unicode/utf8"
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
  ExpectedArgs(args []*Arg) int
  ArgsRequired(args []*Arg) bool
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

func (self *Opts) AppendCommand(op string, i int, spec CommandSpec) {
  self.Commands = append(self.Commands, NewCommand(op, i, spec))
}

func (self *Opts) AppendExtraArg(arg string, i int) {
  self.ExtraArgs = append(self.ExtraArgs, NewArg(arg, i))
}

// Sorts arguments between Commands (beginning with a "-" and a letter for id)
// and ExtraArgs, not tied to any Command. 0 or 1 ExtraArgs are expected, and
// certain Commands expect no sibling Commands.
func (self *Opts) Parse() error {
  nextLiteral := false
  for i, arg := range self.AllArgs {
    if !strings.HasPrefix(arg, "-") || nextLiteral {
      nextLiteral = false
      if lcmd := self.LastCommand(); lcmd != nil && lcmd.ExpectedArgs() > 0 {
        lcmd.AppendArg(arg, i)
      } else {
        self.AppendExtraArg(arg, i)
      }
      continue
    }
    if arg == "-" { // bare "-" escapes the next argument
      nextLiteral = true
      continue
    }
    if arg == "--" { // stop Commands, rest go to extra args
      extrasIndex := i + 1
      for k, arg := range self.AllArgs[extrasIndex:] {
        self.AppendExtraArg(arg, k + extrasIndex)
      }
      break
    }
    var cmd, cmdArg string
    var spec CommandSpec
    prefix := "-"
    if strings.HasPrefix(arg, "--") {
      prefix = "--"
      kv := strings.SplitN(arg[2:], "=", 2)
      cmd = kv[0]
      spec = self.MenuWord(cmd)
      if len(kv) == 2 {
        cmdArg = kv[1]
      }
    } else {
      cmdRune, size := utf8.DecodeRuneInString(arg[1:])
      cmd = string(cmdRune)
      spec = self.MenuLetter(cmdRune)
      cmdArg = arg[1 + size:]
    }
    self.AppendCommand(cmd, i, spec)
    if spec == nil {
      return fmt.Errorf("unknown option: %s in: %s", self.LastCommand(), arg)
    }
    if cmdArg == "" {
      continue
    }
    if lcmd := self.LastCommand(); lcmd.ExpectedArgs() > 0 {
      lcmd.AppendArg(cmdArg, i)
      continue
    }
    if prefix == "--" {
      return fmt.Errorf("--%s takes no values, but saw %s", cmd, arg)
    }
    // series of single-letter flags
    for _, flag := range cmdArg {
      spec = self.MenuLetter(flag)
      self.AppendCommand(string(flag), i, spec)
      if spec == nil {
        return fmt.Errorf("unknown option: %s in: %s", self.LastCommand(), arg)
      }
    }
  }
  for _, opt := range self.Commands {
    if opt.ArgsRequired() {
      return fmt.Errorf("%s is missing a value", opt)
    }
  }
  return nil
}

