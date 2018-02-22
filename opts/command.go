package opts

type Command struct{
  *Arg
  CommandSpec
  Args []*Arg
}

func NewCommand(text string, i int, spec CommandSpec) *Command {
  self := new(Command)
  self.Arg = NewArg(text, i)
  self.CommandSpec = spec
  self.Args = make([]*Arg, 0)
  return self
}

func (self *Command) AppendArg(arg string, i int) {
  self.Args = append(self.Args, NewArg(arg, i))
}

func (self *Command) ExpectedArgs() int {
  return self.CommandSpec.ExpectedArgs(self.Args)
}

func (self *Command) ArgsRequired() bool {
  return self.CommandSpec.ArgsRequired(self.Args)
}

func (self *Command) String() string {
  if len(self.Text) == 1 {
    return "-" + self.Text
  }
  return "--" + self.Text
}

