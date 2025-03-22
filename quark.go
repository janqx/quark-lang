package quark

import (
	"fmt"
	"math"
	"strconv"
)

var (
	Null  Object = &NullObject{}
	True  Object = &BoolObject{Value: true}
	False Object = &BoolObject{Value: false}
)

type CallableFunction func(ctx *Context, args []Object) (Object, error)

const (
	VersionMajor = 0
	VersionMinor = 1
	VersionPatch = 0

	MaxStackSize        = 1024
	MaxCallFrameSize    = 128
	MaxExportObjectSize = 128
	MaxStringLength     = math.MaxInt32
	MaxListLength       = math.MaxInt32
	MaxDictLength       = math.MaxInt32

	MinIntCacheRange = -128
	MaxIntCacheRange = 127

	SourceFileExt = ".qk"
)

func ToBool(x Object) bool {
	return x.ToBool()
}

func ToInt(x Object) (int64, error) {
	switch x := x.(type) {
	case *NullObject:
		return 0, nil
	case *IntObject:
		return x.Value, nil
	case *FloatObject:
		return int64(x.Value), nil
	case *StringObject:
		if result, err := strconv.ParseInt(x.Value, 10, 64); err == nil {
			return result, nil
		}
	}
	return 0, nil
}

func ToFloat(x Object) (float64, error) {
	switch x := x.(type) {
	case *NullObject:
		return 0, nil
	case *IntObject:
		return float64(x.Value), nil
	case *FloatObject:
		return x.Value, nil
	case *StringObject:
		if result, err := strconv.ParseFloat(x.Value, 64); err == nil {
			return result, nil
		}
	}
	return 0, nil
}

func ToString(x Object) (string, error) {
	return x.ToString(), nil
}

func ToInterface(x Object) (interface{}, error) {
	switch x := x.(type) {
	case *NullObject:
		return nil, nil
	case *IntObject:
		return x.Value, nil
	case *FloatObject:
		return x.Value, nil
	case *StringObject:
		return x.Value, nil
	case *ListObject:
		return x.Value, nil
	case *DictObject:
		return x.Value, nil
	}
	return 0, nil
}

func FromInterface(x interface{}) (Object, error) {
	if x == nil {
		return Null, nil
	}
	switch x := x.(type) {
	case Object:
		return x, nil
	case bool:
		if x {
			return True, nil
		}
		return False, nil
	case int:
		return NewInt(int64(x)), nil
	case int64:
		return NewInt(x), nil
	case float32:
		return NewFloat(float64(x)), nil
	case float64:
		return NewFloat(x), nil
	case string:
		return NewString(x), nil
	case []Object:
		return NewList(x), nil
	case map[string]Object:
		return NewDict(x), nil
	}
	return nil, fmt.Errorf("cannot convert to object: %T", x)
}

func FromBool(x bool) Object {
	if x {
		return True
	} else {
		return False
	}
}
