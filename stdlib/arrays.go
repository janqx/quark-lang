package stdlib

import (
	"github.com/janqx/quark-lang/v1"
)

var arraysModule = map[string]quark.Object{
	"createWithLength": quark.NewBuiltinFunction("createWithLength", _createWithLength, 1),
	"fill":             quark.NewBuiltinFunction("fill", _fill, 2),
}

func _createWithLength(ctx *quark.Context, args []quark.Object) (quark.Object, error) {
	value := args[0].(*quark.IntObject)
	list := make([]quark.Object, value.Value)
	for i := 0; i < int(value.Value); i++ {
		list[i] = quark.Null
	}
	return quark.NewList(list), nil
}

func _fill(ctx *quark.Context, args []quark.Object) (quark.Object, error) {
	list := args[0].(*quark.ListObject)
	value := args[1].(*quark.IntObject)
	for index := range list.Value {
		list.Value[index] = value
	}
	return list, nil
}
