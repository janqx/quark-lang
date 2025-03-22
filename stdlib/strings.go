package stdlib

import "github.com/janqx/quark-lang/v1"

var stringsModule = map[string]quark.Object{
	"fromCharCode": quark.NewBuiltinFunction("fromCharCode", _fromCharCode, 1),
}

func _fromCharCode(ctx *quark.Context, args []quark.Object) (quark.Object, error) {
	code, ok := args[0].(*quark.IntObject)
	if !ok {
		return nil, quark.ErrInvalidArgument{
			Name:     "code",
			Expected: "Int",
			Found:    args[0].TypeName(),
		}
	}

	return quark.NewString(string(rune(code.Value))), nil
}
