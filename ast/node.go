package ast

import "github.com/janqx/quark-lang/v1/tokenize"

type Node interface {
	Start() tokenize.Position
	End() tokenize.Position
	String() string
	Accept(visitor Visitor)
}
