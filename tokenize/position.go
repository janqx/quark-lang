package tokenize

type Position struct {
	Filename string
	Offset   int // offset, starting at 0
	Line     int // line, starting at 1
	Column   int // column, starting at 1
}

func (p *Position) IsValid() bool {
	return p.Filename != "" && p.Offset >= 0 && p.Line > 0 && p.Column > 0
}

func (p *Position) String() string {
	return ""
}
