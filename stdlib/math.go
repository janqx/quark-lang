package stdlib

import (
	"math"

	"github.com/janqx/quark-lang/v1"
)

var mathModule = map[string]quark.Object{
	"PI":  quark.NewFloat(math.Pi),
	"E":   quark.NewFloat(math.E),
	"abs": quark.NewBuiltinFunction("abs", quark.TransferAFRF(math.Abs), 1),
	"pow": quark.NewBuiltinFunction("pow", quark.TransferAFFRF(math.Pow), 2),
}
