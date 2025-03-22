package quark

import (
	"fmt"
	"hash/fnv"
	"strconv"
)

type Object interface {
	TypeName() string
	Length() (int, error)
	Callable() bool
	HashCode() int
	Copy() (Object, error)

	ToBool() bool
	ToString() string

	IndexGet(index Object) (Object, error)
	IndexSet(index, value Object) error

	AttributeGet(name string) (Object, error)
	AttributeSet(name string, value Object) error

	UnaryBitNot() (Object, error)
	UnaryNot() (Object, error)
	UnaryPlus() (Object, error)
	UnaryMinus() (Object, error)

	BinaryAdd(x Object) (Object, error)
	BinarySub(x Object) (Object, error)
	BinaryMul(x Object) (Object, error)
	BinaryDiv(x Object) (Object, error)
	BinaryMod(x Object) (Object, error)

	BinaryLt(x Object) (Object, error)
	BinaryLte(x Object) (Object, error)
	BinaryGt(x Object) (Object, error)
	BinaryGte(x Object) (Object, error)
	BinaryEq(x Object) (Object, error)
	BinaryNeq(x Object) (Object, error)

	BinaryBitAnd(x Object) (Object, error)
	BinaryBitOr(x Object) (Object, error)
	BinaryBitXor(x Object) (Object, error)
	BinaryBitLhs(x Object) (Object, error)
	BinaryBitRhs(x Object) (Object, error)
}

type ObjectImpl struct {
}

func (o *ObjectImpl) TypeName() string      { panic(ErrNotImplemented) }
func (o *ObjectImpl) Length() (int, error)  { panic(ErrNotImplemented) }
func (o *ObjectImpl) Callable() bool        { panic(ErrNotImplemented) }
func (o *ObjectImpl) HashCode() int         { panic(ErrNotImplemented) }
func (o *ObjectImpl) Copy() (Object, error) { panic(ErrNotImplemented) }

func (o *ObjectImpl) ToBool() bool     { panic(ErrNotImplemented) }
func (o *ObjectImpl) ToString() string { panic(ErrNotImplemented) }

func (o *ObjectImpl) IndexGet(index Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) IndexSet(index, value Object) error    { panic(ErrNotImplemented) }

func (o *ObjectImpl) AttributeGet(name string) (Object, error)     { panic(ErrNotImplemented) }
func (o *ObjectImpl) AttributeSet(name string, value Object) error { panic(ErrNotImplemented) }

func (o *ObjectImpl) UnaryBitNot() (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) UnaryNot() (Object, error)    { panic(ErrNotImplemented) }
func (o *ObjectImpl) UnaryPlus() (Object, error)   { panic(ErrNotImplemented) }
func (o *ObjectImpl) UnaryMinus() (Object, error)  { panic(ErrNotImplemented) }

func (o *ObjectImpl) BinaryAdd(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinarySub(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryMul(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryDiv(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryMod(x Object) (Object, error) { panic(ErrNotImplemented) }

func (o *ObjectImpl) BinaryLt(x Object) (Object, error)  { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryLte(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryGt(x Object) (Object, error)  { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryGte(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryEq(x Object) (Object, error)  { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryNeq(x Object) (Object, error) { panic(ErrNotImplemented) }

func (o *ObjectImpl) BinaryBitAnd(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryBitOr(x Object) (Object, error)  { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryBitXor(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryBitLhs(x Object) (Object, error) { panic(ErrNotImplemented) }
func (o *ObjectImpl) BinaryBitRhs(x Object) (Object, error) { panic(ErrNotImplemented) }

type NullObject struct {
	ObjectImpl
}

func (o *NullObject) TypeName() string {
	return "Null"
}

func (o *NullObject) ToBool() bool {
	return false
}

func (o *NullObject) ToString() string {
	return "null"
}

func (o *NullObject) Callable() bool {
	return false
}

func (o *NullObject) HashCode() int {
	return 0
}

type BoolObject struct {
	ObjectImpl
	Value bool
}

func (o *BoolObject) TypeName() string {
	return "Bool"
}

func (o *BoolObject) ToBool() bool {
	return o.Value
}

func (o *BoolObject) ToString() string {
	if o.Value {
		return "true"
	}
	return "false"
}
func (o *BoolObject) Callable() bool {
	return false
}

func (o *BoolObject) HashCode() int {
	if o.Value {
		return 1
	}
	return 0
}

func (o *BoolObject) BinaryEq(x Object) (Object, error) {
	value, _ := FromInterface(o == x)
	return value, nil
}

type IntObject struct {
	ObjectImpl
	Value int64
}

func (o *IntObject) TypeName() string {
	return "Int"
}

func (o *IntObject) ToBool() bool {
	return o.Value != 0
}

func (o *IntObject) ToString() string {
	return strconv.FormatInt(o.Value, 10)
}

func (o *IntObject) Copy() (Object, error) {
	return NewInt(o.Value), nil
}

func (o *IntObject) Callable() bool {
	return false
}

func (o *IntObject) HashCode() int {
	return int(o.Value)
}

func (o *IntObject) BinaryAdd(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return NewInt(o.Value + x.Value), nil
	case *FloatObject:
		return NewInt(o.Value + int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for +: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinarySub(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return NewInt(o.Value - x.Value), nil
	case *FloatObject:
		return NewInt(o.Value - int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for -: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryMul(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return NewInt(o.Value * x.Value), nil
	case *FloatObject:
		return NewInt(o.Value * int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for /: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryDiv(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return NewInt(o.Value / x.Value), nil
	case *FloatObject:
		return NewInt(o.Value / int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for *: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryMod(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return NewInt(o.Value % x.Value), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for %%: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryLt(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return FromBool(o.Value < x.Value), nil
	case *FloatObject:
		return FromBool(o.Value < int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for <: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryLte(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return FromBool(o.Value <= x.Value), nil
	case *FloatObject:
		return FromBool(o.Value <= int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for <=: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryGt(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return FromBool(o.Value > x.Value), nil
	case *FloatObject:
		return FromBool(o.Value > int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for >: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryGte(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return FromBool(o.Value >= x.Value), nil
	case *FloatObject:
		return FromBool(o.Value >= int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for >=: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryEq(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return FromBool(o.Value == x.Value), nil
	case *FloatObject:
		return FromBool(o.Value == int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for ==: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *IntObject) BinaryNeq(x Object) (Object, error) {
	switch x := x.(type) {
	case *IntObject:
		return FromBool(o.Value != x.Value), nil
	case *FloatObject:
		return FromBool(o.Value != int64(x.Value)), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for !=: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

var _int_cache map[int64]*IntObject = nil

func NewInt(value int64) *IntObject {
	if _int_cache == nil {
		_int_cache = make(map[int64]*IntObject)
		for i := MinIntCacheRange; i <= MaxIntCacheRange; i++ {
			_int_cache[int64(i)] = &IntObject{
				Value: int64(i),
			}
		}
	}
	if value >= MinIntCacheRange && value <= MaxIntCacheRange {
		return _int_cache[value]
	} else {
		return &IntObject{
			Value: value,
		}
	}
}

type FloatObject struct {
	ObjectImpl
	Value float64
}

func (o *FloatObject) ToString() string {
	return fmt.Sprintf("%.12f", o.Value)
}

func NewFloat(value float64) *FloatObject {
	return &FloatObject{
		Value: value,
	}
}

type StringObject struct {
	ObjectImpl
	Value string
}

func (o *StringObject) TypeName() string {
	return "String"
}

func (o *StringObject) ToString() string {
	return o.Value
}

func (o *StringObject) Copy() (Object, error) {
	return NewString(o.Value), nil
}

func (o *StringObject) Length() (int, error) {
	return len(o.Value), nil
}

func (o *StringObject) Callable() bool {
	return false
}

func (o *StringObject) HashCode() int {
	h := fnv.New32a()
	h.Write([]byte(o.Value))
	return int(h.Sum32())
}

func (o *StringObject) BinaryAdd(x Object) (Object, error) {
	switch x := x.(type) {
	case *StringObject:
		return NewString(o.Value + x.Value), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for +: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *StringObject) BinaryEq(x Object) (Object, error) {
	switch x := x.(type) {
	case *StringObject:
		return FromBool(o.Value == x.Value), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for ==: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *StringObject) BinaryNeq(x Object) (Object, error) {
	switch x := x.(type) {
	case *StringObject:
		return FromBool(o.Value != x.Value), nil
	default:
		return nil, fmt.Errorf("unsupported operand type(s) for !=: '%s' and '%s'", o.TypeName(), x.TypeName())
	}
}

func (o *StringObject) IndexGet(index Object) (Object, error) {
	if index, ok := index.(*IntObject); !ok {
		return nil, ErrInvalidIndexType
	} else if index.Value < 0 || int(index.Value) >= len(o.Value) {
		return nil, ErrIndexOutOfRange
	} else {
		return NewString(string(o.Value[index.Value])), nil
	}
}

func NewString(value string) *StringObject {
	return &StringObject{
		Value: value,
	}
}

type ListObject struct {
	ObjectImpl
	Value []Object
}

func (o *ListObject) TypeName() string {
	return "List"
}

func (o *ListObject) ToString() string {
	result := "["
	for i, obj := range o.Value {
		result += obj.ToString()
		if i < len(o.Value)-1 {
			result += ", "
		}
	}
	return result + "]"
}

func (o *ListObject) Length() (int, error) {
	return len(o.Value), nil
}

func (o *ListObject) ToBool() bool {
	return len(o.Value) != 0
}

func (o *ListObject) Copy() (Object, error) {
	return NewList(o.Value), nil
}

func (o *ListObject) Callable() bool {
	return false
}

func (o *ListObject) IndexGet(index Object) (Object, error) {
	if index, ok := index.(*IntObject); !ok {
		return nil, ErrInvalidIndexType
	} else if index.Value < 0 || int(index.Value) >= len(o.Value) {
		return nil, ErrIndexOutOfRange
	} else {
		return o.Value[index.Value], nil
	}
}

func (o *ListObject) IndexSet(index, value Object) error {
	if index, ok := index.(*IntObject); !ok {
		return ErrInvalidIndexType
	} else if index.Value < 0 || int(index.Value) >= len(o.Value) {
		return ErrIndexOutOfRange
	} else {
		o.Value[index.Value] = value
		return nil
	}
}

func NewList(value []Object) *ListObject {
	return &ListObject{
		Value: value,
	}
}

type DictObject struct {
	ObjectImpl
	Value map[string]Object
}

func (o *DictObject) TypeName() string {
	return "Dict"
}

func (o *DictObject) ToString() string {
	if len(o.Value) == 0 {
		return "{}"
	}
	result := "{ "
	index := 0
	for key, value := range o.Value {
		result += key + ": "
		result += value.ToString()
		if index < len(o.Value)-1 {
			result += ", "
		}
		index++
	}
	return result + " }"
}

func (o *DictObject) Length() (int, error) {
	return len(o.Value), nil
}

func (o *DictObject) ToBool() bool {
	return len(o.Value) != 0
}

func (o *DictObject) IndexGet(index Object) (Object, error) {
	key, err := ToString(index)
	if err != nil {
		return nil, err
	}
	if value, ok := o.Value[key]; ok {
		return value, nil
	} else {
		return Null, nil
	}
}

func (o *DictObject) IndexSet(index, value Object) error {
	key, err := ToString(index)
	if err != nil {
		return nil
	}
	o.Value[key] = value
	return nil
}

func (o *DictObject) AttributeGet(name string) (Object, error) {
	if name == "" {
		return nil, ErrInvalidAttributeName
	}
	if value, ok := o.Value[name]; ok {
		return value, nil
	} else {
		return Null, nil
	}
}

func (o *DictObject) AttributeSet(name string, value Object) error {
	if name == "" {
		return ErrInvalidAttributeName
	}
	o.Value[name] = value
	return nil
}

func NewDict(value map[string]Object) *DictObject {
	return &DictObject{
		Value: value,
	}
}

type BuiltinFunctionObject struct {
	ObjectImpl
	Name          string
	Fn            CallableFunction
	NumParameters int
}

func (o *BuiltinFunctionObject) TypeName() string {
	return "BuiltinFunction"
}

func (o *BuiltinFunctionObject) ToString() string {
	return fmt.Sprintf("<builtin-function %s>", o.Name)
}

func (o *BuiltinFunctionObject) Copy() (Object, error) {
	return NewBuiltinFunction(o.Name, o.Fn, o.NumParameters), nil
}

func (o *BuiltinFunctionObject) Callable() bool {
	return true
}

func NewBuiltinFunction(name string, fn CallableFunction, numParameters int) *BuiltinFunctionObject {
	return &BuiltinFunctionObject{
		Name:          name,
		Fn:            fn,
		NumParameters: numParameters,
	}
}

type CompiledFunctionObject struct {
	ObjectImpl
	Name           string
	Instructions   []Instruction
	ParameterNames []string
	SymbolTable    *SymbolTable
}

func (o *CompiledFunctionObject) TypeName() string {
	return "BuiltinFunction"
}

func (o *CompiledFunctionObject) ToString() string {
	return fmt.Sprintf("<builtin-function %s>", o.Name)
}

func (o *CompiledFunctionObject) Copy() (Object, error) {
	return &CompiledFunctionObject{
		Name:           o.Name,
		Instructions:   o.Instructions,
		ParameterNames: o.ParameterNames,
		SymbolTable:    o.SymbolTable,
	}, nil
}

func (o *CompiledFunctionObject) Callable() bool {
	return true
}

type ClosureObject struct {
	ObjectImpl
	Fn     *CompiledFunctionObject
	Outers []Object
}

func (o *ClosureObject) TypeName() string {
	return "Closure"
}

func (o *ClosureObject) ToString() string {
	return fmt.Sprintf("<closure %s>", o.Fn.Name)
}

func (o *ClosureObject) Copy() (Object, error) {
	return &ClosureObject{
		Fn:     o.Fn,
		Outers: o.Outers,
	}, nil
}

func (o *ClosureObject) Callable() bool {
	return true
}

type ObjectRef struct {
	ObjectImpl
	Value Object
}
