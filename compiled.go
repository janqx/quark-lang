package quark

import "fmt"

type compiled struct {
	entryFunction     *CompiledFunctionObject
	compiledFunctions []*CompiledFunctionObject
}

func (c *compiled) PrintInstructionList() {
	fmt.Printf("PrintInstructionList(%d):\n", len(c.entryFunction.Instructions))
	for _, fn := range c.compiledFunctions {
		fmt.Println("---> compiled-function: " + fn.Name)
		var index int = 0
		for _, ins := range fn.Instructions {
			fmt.Printf("#%d: %s\n", index, ins.String())
			index++
		}
	}
	fmt.Println()
}
