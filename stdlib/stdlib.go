package stdlib

import "github.com/janqx/quark-lang/v1"

var modules = map[string]map[string]quark.Object{
	"fmt":     nil,
	"os":      nil,
	"math":    mathModule,
	"time":    nil,
	"strings": stringsModule,
	"arrays":  arraysModule,
}

func LoadModules() map[string]map[string]quark.Object {
	return modules
}
