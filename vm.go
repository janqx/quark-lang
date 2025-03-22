package quark

import (
	"fmt"
	"sync/atomic"
)

type CallFrame struct {
	fn     *CompiledFunctionObject
	outers []Object
	ip     int
	bp     int
}

type VM struct {
	ctx *Context
}

func NewVM(ctx *Context) *VM {
	return &VM{
		ctx: ctx,
	}
}

// 为了接下来执行的callable对象做准备
func (vm *VM) Prepare(callee Object, argc int) error {
	return vm.call(callee, argc)
}

// 执行callable对象，并返回函数返回值
func (vm *VM) Execute() (Object, error) {
	if err := vm.execute(); err != nil {
		return nil, err
	}
	return vm.pop(), nil
}

func (vm *VM) execute() error {
	ctx := vm.ctx
	for ctx.ip+1 < len(ctx.currentFrame.fn.Instructions) && atomic.LoadInt32(&(ctx.abortFlag)) == 0 {
		ctx.ip++
		inst := ctx.currentFrame.fn.Instructions[ctx.ip]
		switch inst.Opcode() {
		case OpNop:
		case OpLoadNull:
			vm.push(Null)
		case OpLoadTrue:
			vm.push(True)
		case OpLoadFalse:
			vm.push(False)
		case OpLoadConst:
			vm.push(ctx.constants[inst.Operand()])
		case OpLoadLocal:
			vm.push(vm.getLocal(int(inst.Operand())))
		case OpLoadOuter:
			vm.push(vm.getOuter(int(inst.Operand())))
		case OpLoadGlobal:
			vm.push(vm.getGlobal(int(inst.Operand())))
		case OpLoadIndex:
			index := vm.pop()
			obj := vm.pop()
			value, err := obj.IndexGet(index)
			if err != nil {
				return err
			}
			vm.push(value)
		case OpLoadAttribute:
			attr := vm.pop()
			obj := vm.pop()
			value, err := obj.AttributeGet(attr.(*StringObject).Value)
			if err != nil {
				return err
			}
			vm.push(value)
		case OpStoreLocal:
			vm.setLocal(int(inst.Operand()), vm.pop())
		case OpStoreOuter:
			vm.setOuter(int(inst.Operand()), vm.pop())
		case OpStoreGlobal:
			vm.setGlobal(int(inst.Operand()), vm.pop())
		case OpStoreIndex:
			index := vm.pop()
			list := vm.pop()
			value := vm.pop()
			if err := list.IndexSet(index, value); err != nil {
				return err
			}
		case OpStoreAttribute:
			attr := vm.pop()
			obj := vm.pop()
			value := vm.pop()
			if err := obj.AttributeSet(attr.(*StringObject).Value, value); err != nil {
				return err
			}
		case OpBuildList:
			if err := vm.buildList(int(inst.Operand())); err != nil {
				return err
			}
		case OpBuildDict:
			if err := vm.buildDict(int(inst.Operand())); err != nil {
				return err
			}
		case OpUnaryBitNot, OpUnaryNot, OpUnaryPlus, OpUnaryMinus:
			obj, err := vm.unaryOp(inst.Opcode(), vm.pop())
			if err != nil {
				return err
			}
			vm.push(obj)
		case OpBinaryAdd,
			OpBinarySub,
			OpBinaryMul,
			OpBinaryDiv,
			OpBinaryMod,
			OpBinaryLT,
			OpBinaryLTE,
			OpBinaryGT,
			OpBinaryGTE,
			OpBinaryEQ,
			OpBinaryNEQ,
			OpBinaryBitAnd,
			OpBinaryBitOr,
			OpBinaryBitXor,
			OpBinaryBitLhs,
			OpBinaryBitRhs:
			ctx.sp -= 2
			obj, err := vm.binaryOp(inst.Opcode(), ctx.stack[ctx.sp], ctx.stack[ctx.sp+1])
			if err != nil {
				return err
			}
			vm.push(obj)
		case OpJump:
			ctx.ip = int(inst.Operand()) - 1
		case OpJumpIfFalse:
			if !vm.pop().ToBool() {
				ctx.ip = int(inst.Operand()) - 1
			}
		case OpJumpIfFalseOrPop:
			if !vm.peek().ToBool() {
				ctx.ip = int(inst.Operand()) - 1
			} else {
				vm.pop()
			}
		case OpJumpIfTrueOrPop:
			if vm.peek().ToBool() {
				ctx.ip = int(inst.Operand()) - 1
			} else {
				vm.pop()
			}
		case OpClosure:
			closure, err := vm.makeClosure(vm.pop())
			if err != nil {
				return err
			}
			vm.push(closure)
		case OpCall:
			if err := vm.call(vm.pop(), int(inst.Operand())); err != nil {
				return err
			}
		case OpReturn:
			result := vm.pop()
			ctx.ip = ctx.currentFrame.ip
			ctx.sp = ctx.currentFrame.bp
			vm.push(result)
			ctx.fp--
			ctx.currentFrame = ctx.frames[ctx.fp]
		case OpRemoveTop:
			vm.pop()
		case OpImport:
			// modulePath := ctx.constants[inst.Operand()].(*StringObject).Value
			// moduleAbsolute, err := filepath.Abs(filepath.Join(ctx.ImportBasePath, modulePath))
			// if err != nil {
			// 	return err
			// }
			// if module, ok := ctx.compiledModules[moduleAbsolute]; ok {
			// 	vm.push(module)
			// } else if module, ok := ctx.builtinModules[modulePath]; ok {
			// 	vm.push(module)
			// } else {
			// source, err := os.ReadFile(moduleAbsolute)
			// if err != nil {
			// 	return err
			// }
			// p := parser.NewParser(moduleAbsolute, source)
			// chunk, err := p.Parse()
			// if err != nil {
			// 	return err
			// }
			// compiler := NewCompiler(ctx)
			// compiled := compiler.Compile(chunk)
			// if compiler.Err != nil {
			// 	return err
			// }
			// compiled.PrintInstructionList()
			// _, err = NewVM(ctx).ExecuteWithCompiledFunction(compiled.EntryFunction)
			// if err != nil {
			// 	return err
			// }
			// ctx.exportObjectIndex--
			// exportObject := ctx.exportObjects[ctx.exportObjectIndex]
			// ctx.compiledModules[moduleAbsolute] = exportObject
			// vm.push(exportObject)
			// }
		case OpExport:
			ctx.exportObjects[ctx.exportObjectIndex] = vm.pop()
			ctx.exportObjectIndex++

			result := Null
			ctx.ip = ctx.currentFrame.ip
			ctx.sp = ctx.currentFrame.bp
			vm.push(result)
			ctx.fp--
			ctx.currentFrame = ctx.frames[ctx.fp]
			return nil
		case OpDebugger:
			fmt.Println("debugger")
		default:
			return ErrInvalidOpcode
		}
	}
	return nil
}

func (vm *VM) unaryOp(opcode Opcode, x Object) (Object, error) {
	switch opcode {
	case OpUnaryBitNot:
		return x.UnaryBitNot()
	case OpUnaryNot:
		return x.UnaryNot()
	case OpUnaryPlus:
		return x.UnaryPlus()
	case OpUnaryMinus:
		return x.UnaryMinus()
	default:
		return nil, ErrInvalidOpcode
	}
}

func (vm *VM) binaryOp(opcode Opcode, left, right Object) (Object, error) {
	switch opcode {
	case OpBinaryAdd:
		return left.BinaryAdd(right)
	case OpBinarySub:
		return left.BinarySub(right)
	case OpBinaryMul:
		return left.BinaryMul(right)
	case OpBinaryDiv:
		return left.BinaryDiv(right)
	case OpBinaryMod:
		return left.BinaryMod(right)
	case OpBinaryLT:
		return left.BinaryLt(right)
	case OpBinaryLTE:
		return left.BinaryLte(right)
	case OpBinaryGT:
		return left.BinaryGt(right)
	case OpBinaryGTE:
		return left.BinaryGte(right)
	case OpBinaryEQ:
		return left.BinaryEq(right)
	case OpBinaryNEQ:
		return left.BinaryNeq(right)
	case OpBinaryBitAnd:
		return left.BinaryBitAnd(right)
	case OpBinaryBitOr:
		return left.BinaryBitOr(right)
	case OpBinaryBitXor:
		return left.BinaryBitXor(right)
	case OpBinaryBitLhs:
		return left.BinaryBitLhs(right)
	case OpBinaryBitRhs:
		return left.BinaryBitRhs(right)
	default:
		return nil, ErrInvalidOpcode
	}
}

func (vm *VM) buildList(count int) error {
	list := make([]Object, count)
	for i := 0; i < count; i++ {
		list[i] = vm.ctx.stack[vm.ctx.sp-count+i]
	}
	vm.ctx.sp -= count
	vm.push(&ListObject{
		Value: list,
	})
	return nil
}

func (vm *VM) buildDict(count int) error {
	m := make(map[string]Object)
	for i := 0; i < count; i++ {
		value := vm.pop()
		key := vm.pop()
		m[key.(*StringObject).Value] = value
	}
	vm.push(&DictObject{
		Value: m,
	})
	return nil
}

func (vm *VM) call(callee Object, argc int) error {
	if !callee.Callable() {
		return fmt.Errorf("can't call object: %s", callee.TypeName())
	}

	var args []Object = nil

	if argc > 0 {
		args = make([]Object, argc)
		for i := 0; i < argc; i++ {
			args[i] = vm.ctx.stack[vm.ctx.sp-argc+i]
		}
		vm.ctx.sp -= argc
	}

	switch callee := callee.(type) {
	case *CompiledFunctionObject:
		return vm.callCompiledFunction(callee, args)
	case *ClosureObject:
		return vm.callClosure(callee, args)
	case *BuiltinFunctionObject:
		return vm.callBuiltinFunction(callee, args)
	default:
		return fmt.Errorf("not implement call type: %s", callee.TypeName())
	}
}

func (vm *VM) makeClosure(fn Object) (*ClosureObject, error) {
	compiledFn, ok := fn.(*CompiledFunctionObject)
	if !ok {
		return nil, fmt.Errorf("is not a compiled-function: %s", fn.TypeName())
	}
	numOuters := compiledFn.SymbolTable.OuterCount
	outers := make([]Object, numOuters)
	if numOuters > 0 {
		for _, symbol := range compiledFn.SymbolTable.Symbols {
			if symbol.Scope == ScopeOuter {
				if symbol.OuterScope == ScopeLocal {
					switch obj := vm.ctx.stack[vm.ctx.currentFrame.bp+symbol.OuterIndex].(type) {
					case *ObjectRef:
						outers[symbol.Index] = obj
					default:
						outers[symbol.Index] = &ObjectRef{
							Value: obj,
						}
						vm.ctx.stack[vm.ctx.currentFrame.bp+symbol.OuterIndex] = outers[symbol.Index]
					}
				} else if symbol.OuterScope == ScopeOuter {
					switch obj := vm.ctx.currentFrame.outers[symbol.OuterIndex].(type) {
					case *ObjectRef:
						outers[symbol.Index] = obj
					default:
						outers[symbol.Index] = &ObjectRef{
							Value: obj,
						}
						vm.ctx.currentFrame.outers[symbol.OuterIndex] = outers[symbol.Index]
					}
				}
			}
		}
	}
	return &ClosureObject{
		Fn:     compiledFn,
		Outers: outers,
	}, nil
}

func (vm *VM) callCompiledFunction(fn *CompiledFunctionObject, args []Object) error {
	closure, err := vm.makeClosure(fn)
	if err != nil {
		return err
	}
	return vm.callClosure(closure, args)
}

func (vm *VM) callClosure(closure *ClosureObject, args []Object) error {
	frame := &CallFrame{
		fn:     closure.Fn,
		outers: closure.Outers,
		ip:     vm.ctx.ip,
		bp:     vm.ctx.sp,
	}

	vm.ctx.sp += closure.Fn.SymbolTable.LocalCount

	if vm.ctx.sp >= MaxStackSize {
		return ErrStackOverflow
	}

	if len(args) > 0 {
		copy(vm.ctx.stack[frame.bp:frame.bp+len(args)], args)
	}

	vm.ctx.fp++
	vm.ctx.frames[vm.ctx.fp] = frame
	vm.ctx.currentFrame = frame
	vm.ctx.ip = -1

	return nil
}

func (vm *VM) callBuiltinFunction(fn *BuiltinFunctionObject, args []Object) error {
	if fn.NumParameters != len(args) {
		return ErrWrongNumberArguments
	}
	if result, err := fn.Fn(vm.ctx, args); err != nil {
		return err
	} else {
		vm.push(result)
		return nil
	}
}

func (vm *VM) push(obj Object) {
	vm.ctx.stack[vm.ctx.sp] = obj
	vm.ctx.sp++
}

func (vm *VM) pop() Object {
	vm.ctx.sp--
	return vm.ctx.stack[vm.ctx.sp]
}

func (vm *VM) peek() Object {
	return vm.ctx.stack[vm.ctx.sp-1]
}

func (vm *VM) getLocal(index int) Object {
	value := vm.ctx.stack[vm.ctx.currentFrame.bp+index]
	if ref, ok := value.(*ObjectRef); ok {
		return ref.Value
	} else {
		return value
	}
}

func (vm *VM) setLocal(index int, value Object) {
	if ref, ok := vm.ctx.stack[vm.ctx.currentFrame.bp+index].(*ObjectRef); ok {
		ref.Value = value
	} else {
		vm.ctx.stack[vm.ctx.currentFrame.bp+index] = value
	}
}

func (vm *VM) getOuter(index int) Object {
	value := vm.ctx.currentFrame.outers[index]
	if ref, ok := value.(*ObjectRef); ok {
		return ref.Value
	} else {
		return value
	}
}

func (vm *VM) setOuter(index int, value Object) {
	if ref, ok := vm.ctx.currentFrame.outers[index].(*ObjectRef); ok {
		ref.Value = value
	} else {
		vm.ctx.currentFrame.outers[index] = value
	}
}

func (vm *VM) getGlobal(index int) Object {
	value := vm.ctx.globals[index]
	if ref, ok := value.(*ObjectRef); ok {
		return ref.Value
	} else {
		return value
	}
}

func (vm *VM) setGlobal(index int, value Object) {
	if ref, ok := vm.ctx.globals[index].(*ObjectRef); ok {
		ref.Value = value
	} else {
		vm.ctx.globals[index] = value
	}
}
