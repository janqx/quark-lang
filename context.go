package quark

import (
	"fmt"
	"os"
)

type InterpreterMode uint8

const (
	ModeNormal InterpreterMode = iota
	ModeREPL
)

type Context struct {
	Mode           InterpreterMode
	AllowImport    bool
	ImportBasePath string

	// used for vm
	constants         []Object
	globals           []Object
	builtinModules    map[string]Object
	compiledModules   map[string]Object
	exportObjects     [MaxExportObjectSize]Object
	exportObjectIndex int
	stack             [MaxStackSize]Object
	sp                int
	frames            [MaxCallFrameSize]*CallFrame
	fp                int
	currentFrame      *CallFrame
	ip                int
	abortFlag         int32

	// used for compiler
	globalSymbolTable *SymbolTable
	intConstantMap    map[int64]int
	floatConstantMap  map[float64]int
	stringConstantMap map[string]int

	err error
}

func NewContext(mode InterpreterMode, stdlibModules map[string]map[string]Object) *Context {
	ctx := &Context{}
	ctx.Mode = mode
	ctx.AllowImport = true
	ctx.ImportBasePath, _ = os.Getwd()

	ctx.constants = make([]Object, 0)
	ctx.globals = make([]Object, 0)
	ctx.builtinModules = make(map[string]Object)
	ctx.compiledModules = make(map[string]Object)

	ctx.globalSymbolTable = NewSymbolTable(nil, TypeFunction)
	ctx.intConstantMap = make(map[int64]int)
	ctx.floatConstantMap = make(map[float64]int)
	ctx.stringConstantMap = make(map[string]int)

	topFn := &CompiledFunctionObject{
		Name:           "<top-function>",
		Instructions:   []Instruction{},
		ParameterNames: []string{},
		SymbolTable:    ctx.globalSymbolTable,
	}

	topFrame := &CallFrame{
		fn:     topFn,
		outers: []Object{},
		ip:     -1,
		bp:     0,
	}

	ctx.sp = 0
	ctx.frames[0] = topFrame
	ctx.fp = 0
	ctx.currentFrame = topFrame
	ctx.ip = -1
	ctx.abortFlag = 0

	// used for repl
	ctx.setGlobal("__REPL_RESULT_VALUE__", Null)

	// initialize builtin objects
	for name, value := range builtinObjects {
		ctx.setGlobal(name, value)
	}

	// initialize stdlib modules
	for name, value := range stdlibModules {
		ctx.builtinModules[name] = &DictObject{
			Value: value,
		}
	}

	return ctx
}

func (c *Context) appendConstant(value Object) int {
	c.constants = append(c.constants, value)
	return len(c.constants) - 1
}

func (c *Context) addIntConstant(value int64) int {
	if index, ok := c.intConstantMap[value]; ok {
		return index
	}
	index := c.appendConstant(NewInt(value))
	c.intConstantMap[value] = index
	return index
}

func (c *Context) addFloatConstant(value float64) int {
	if index, ok := c.floatConstantMap[value]; ok {
		return index
	}
	index := c.appendConstant(NewFloat(value))
	c.floatConstantMap[value] = index
	return index
}

func (c *Context) addStringConstant(value string) int {
	if index, ok := c.stringConstantMap[value]; ok {
		return index
	}
	index := c.appendConstant(NewString(value))
	c.stringConstantMap[value] = index
	return index
}

func (c *Context) addGlobalSymbol(name string) *Symbol {
	symbol := c.globalSymbolTable.AddGlobalSymbol(name)
	if len(c.globals) <= symbol.Index {
		newGlobals := make([]Object, symbol.Index+1)
		copy(newGlobals, c.globals[:])
		c.globals = newGlobals
	}
	return symbol
}

func (c *Context) setGlobal(name string, value Object) {
	symbol := c.addGlobalSymbol(name)
	c.globals[symbol.Index] = value
}

func (c *Context) getGlobal(name string) Object {
	symbol := c.globalSymbolTable.FindSymbol(name)
	if symbol.Scope == ScopeGlobal {
		return c.globals[symbol.Index]
	} else {
		return nil
	}
}

func (c *Context) GetStackTraceback() string {
	return ""
}

func (c *Context) ThrowErrorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func (c *Context) ThrowErrorMessage(message ErrorMessage) {
	c.ThrowErrorf("%s:%d: %s", message.Filename, message.Line, message.Message)
}
