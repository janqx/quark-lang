package quark

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/janqx/quark-lang/v1/parser"
)

var builtinObjects = map[string]Object{
	"print":     NewBuiltinFunction("print", _print, 1),
	"println":   NewBuiltinFunction("println", _println, 1),
	"panic":     NewBuiltinFunction("panic", _panic, 1),
	"input":     NewBuiltinFunction("input", _input, 1),
	"format":    nil,
	"copy":      nil,
	"length":    NewBuiltinFunction("length", _length, 1),
	"typename":  NewBuiltinFunction("typename", _typename, 1),
	"typeid":    nil,
	"append":    nil,
	"delete":    nil,
	"import":    NewBuiltinFunction("import", _import, 1),
	"to_bool":   NewBuiltinFunction("to_bool", _to_bool, 1),
	"to_int":    NewBuiltinFunction("to_int", _to_int, 1),
	"to_float":  NewBuiltinFunction("to_float", _to_float, 1),
	"to_string": NewBuiltinFunction("to_string", _to_string, 1),
	"chr":       NewBuiltinFunction("chr", _chr, 1),
}

func _print(ctx *Context, args []Object) (Object, error) {
	ss := make([]string, 0)
	for _, arg := range args {
		s, _ := ToString(arg)
		ss = append(ss, s)
	}
	fmt.Print(strings.Join(ss, " "))
	return Null, nil
}

func _println(ctx *Context, args []Object) (Object, error) {
	result, err := _print(ctx, args)
	fmt.Println()
	return result, err
}

func _panic(ctx *Context, args []Object) (Object, error) {
	fmt.Println("PANIC:")
	result, err := _print(ctx, args)
	fmt.Println()
	return result, err
}

func _input(ctx *Context, args []Object) (Object, error) {
	var input string
	prompt, _ := ToString(args[0])
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return FromInterface(input)
}

func _length(ctx *Context, args []Object) (Object, error) {
	length, err := args[0].Length()
	if err != nil {
		return nil, err
	}
	return FromInterface(length)
}

func _typename(ctx *Context, args []Object) (Object, error) {
	return FromInterface(args[0].TypeName())
}

func _import(ctx *Context, args []Object) (Object, error) {
	filename, ok := args[0].(*StringObject)
	if !ok {
		return nil, ErrInvalidArgument{
			Name:     "module",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	modulePath := filename.Value

	var err error
	var moduleAbsolute string
	if moduleAbsolute, err = filepath.Abs(filepath.Join(ctx.ImportBasePath, modulePath)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	if module, ok := ctx.compiledModules[moduleAbsolute]; ok {
		return module, nil
	} else if module, ok := ctx.builtinModules[modulePath]; ok {
		return module, nil
	}

	source, err := os.ReadFile(moduleAbsolute)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	p := parser.NewParser(filename.Value, source)
	chunk, err := p.Parse()
	if err != nil {
		return nil, err
	}

	compiler := NewCompiler(ctx, nil)
	compiled, err := compiler.Compile(chunk)
	if err != nil {
		return nil, err
	}

	compiled.PrintInstructionList()

	vm := NewVM(ctx)
	err = vm.Prepare(compiled.entryFunction, 0)
	if err != nil {
		return nil, err
	}

	ctx.exportObjectIndex--
	module := ctx.exportObjects[ctx.exportObjectIndex]

	ctx.compiledModules[moduleAbsolute] = module

	return module, nil
}

func _to_bool(ctx *Context, args []Object) (Object, error) {
	return FromInterface(ToBool(args[0]))
}

func _to_int(ctx *Context, args []Object) (Object, error) {
	if value, err := ToInt(args[0]); err != nil {
		return nil, err
	} else {
		return FromInterface(value)
	}
}

func _to_float(ctx *Context, args []Object) (Object, error) {
	if value, err := ToFloat(args[0]); err != nil {
		return nil, err
	} else {
		return FromInterface(value)
	}
}

func _to_string(ctx *Context, args []Object) (Object, error) {
	if value, err := ToString(args[0]); err != nil {
		return nil, err
	} else {
		return FromInterface(value)
	}
}

func _chr(ctx *Context, args []Object) (Object, error) {
	value, ok := args[0].(*IntObject)
	if !ok {
		return nil, ErrInvalidArgument{
			Name:     "code",
			Expected: "Int",
			Found:    args[0].TypeName(),
		}
	}

	return NewString(string(rune(value.Value))), nil
}
