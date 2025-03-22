package quark

type Bytecode struct {
	MainFunction *CompiledFunctionObject
	Constants    []Object
}
