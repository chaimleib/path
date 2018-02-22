package opts

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

