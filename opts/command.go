package opts

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

func (self *Command) AppendArg(arg string, i int) {
  self.Args = append(self.Args, NewArg(arg, i))
}

