package quark

import (
	"errors"
	"fmt"
)

var (
	ErrStackOverflow          = errors.New("stack overflow")
	ErrNotImplemented         = errors.New("not implemented")
	ErrWrongNumberArguments   = errors.New("wrong number of arguments")
	ErrInvalidOpcode          = errors.New("invalid opcode")
	ErrInvalidIndexType       = errors.New("invalid index type")
	ErrIndexOutOfRange        = errors.New("index out of range")
	ErrInvalidAttributeName   = errors.New("invalid attribute name")
	ErrUndeclaredSymbolName   = errors.New("undeclared name")
	ErrTypeIsNotSubscriptable = errors.New("type is not subscriptable")
	ErrInvalidOperator        = errors.New("invalid operator")
	ErrInvalidModuleName      = errors.New("invalid module name")
	ErrNotFoundModule         = errors.New("not found module")
)

type ErrorMessage struct {
	Filename string
	Offset   int
	Line     int
	Column   int
	Message  string
}

type ErrInvalidArgument struct {
	Name     string
	Expected string
	Found    string
}

func NewErrInvalidArgument(name string, expected string, found string) ErrInvalidArgument {
	return ErrInvalidArgument{
		Name:     name,
		Expected: expected,
		Found:    found,
	}
}

func (e ErrInvalidArgument) Error() string {
	return fmt.Sprintf("invalid type for argument '%s': expected %s, found %s", e.Name, e.Expected, e.Found)
}
