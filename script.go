package quark

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/janqx/quark-lang/v1/parser"
)

type Script struct {
	ctx       *Context
	variables map[string]*Variable
}

func NewScript(ctx *Context) *Script {
	return &Script{
		ctx:       ctx,
		variables: make(map[string]*Variable),
	}
}

func (s *Script) SetContext(ctx *Context) {
	s.ctx = ctx
}

func (s *Script) AddVariable(name string, value interface{}) error {
	variable, err := NewVariable(name, value)
	if err != nil {
		return err
	}
	s.variables[name] = variable
	return nil
}

func (s *Script) GetVariable(name string) (interface{}, bool, error) {
	if variable, ok := s.variables[name]; ok {
		object, err := ToInterface(variable.Value())
		return object, err == nil, err
	}
	return nil, false, nil
}

func (s *Script) RemoveVariable(name string) {
	delete(s.variables, name)
}

func (s *Script) ExistsVariable(name string) bool {
	_, exists := s.variables[name]
	return exists
}

func (s *Script) RunString(source string) (Object, error) {
	compiled, err := s.compile("<repl>", source)
	if err != nil {
		return nil, err
	}
	compiled.PrintInstructionList()
	vm := NewVM(s.ctx)
	err = vm.Prepare(compiled.entryFunction, 0)
	if err != nil {
		return nil, err
	}
	return vm.Execute()
}

func (s *Script) RunFile(filename string) error {
	var err error
	var fullpath string
	if ext := filepath.Ext(filename); ext != SourceFileExt {
		return fmt.Errorf("invalid ext name: %s, except: %s", ext, SourceFileExt)
	}
	if fullpath, err = filepath.Abs(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	source, err := os.ReadFile(fullpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	compiled, err := s.compile(fullpath, string(source))
	if err != nil {
		return err
	}
	vm := NewVM(s.ctx)
	err = vm.Prepare(compiled.entryFunction, 0)
	if err != nil {
		return err
	}
	_, err = vm.Execute()
	return err
}

func (s *Script) compile(filename string, source string) (*compiled, error) {
	p := parser.NewParser(filename, []byte(source))
	chunk, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return NewCompiler(s.ctx, nil).Compile(chunk)
}
