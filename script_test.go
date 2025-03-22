package quark_test

import (
	"testing"

	"github.com/janqx/quark-lang/v1"
	"github.com/janqx/quark-lang/v1/stdlib"
)

func TestScript_RunString(t *testing.T) {
	ctx := quark.NewContext(quark.ModeNormal, stdlib.LoadModules())
	script := quark.NewScript(ctx)
	script.RunString(`

fn swap(list, i, j) {
  list[i], list[j] = list[j], list[i]
}

fn quickSort(list, left, right) {
  if left >= right {
    return
  }
  i, j, pivot = left, right, left
  for ;i < j; {
    for ;i < j && list[j] >= list[pivot]; {
      j = j - 1
    }
    for ;i < j && list[i] <= list[pivot]; {
      i = i + 1
    }
    swap(list, i, j)
  }
  swap(list, i, pivot)
  quickSort(list, left, i - 1)
  quickSort(list, i + 1, right)
}

list = [3, 5, 1, 7, 9, 2, 6, 4]
quickSort(list, 0, length(list)-1)
print(list)


	`)
}

func TestScript_RunString_Import(t *testing.T) {
	ctx := quark.NewContext(quark.ModeNormal, stdlib.LoadModules())
	script := quark.NewScript(ctx)
	script.RunString(`

  fn fib(n) {
    if n < 3 {
      return 1
    }
    return fib(n-1)+fib(n-2)
  }

  println(fib(35))

	`)
}

func TestScript_RunFile_Brainfuck(t *testing.T) {
	ctx := quark.NewContext(quark.ModeNormal, stdlib.LoadModules())
	script := quark.NewScript(ctx)
	err := script.RunFile(`E:\jqx\projects\quark-lang\example\brainfuck.qk`)
	if err != nil {
		panic(err)
	}
}
