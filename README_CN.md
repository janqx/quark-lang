# 娜贝脚本语言

***基于GO开发的一种微小、动态、快速的的脚本语言.***

![](images/1.jpeg)

## Features
- REPL
- 模块
- 易嵌入GoLang
- 代码易读

## 使用
```
go get github.com/janqx/narbe/v1
```

## 基本例子
```go
// go run narbe/cli/main.go ./basic.nb

fmt := import("fmt")

fn each(seq, fn) {
    for i := 0; i < len(seq); i = i + 1 {
        fn(seq[i])
    }
}

fn sum(init, seq) {
    each(seq, fn(x) {
        init += x
    })
    return init
}

fmt.println(sum(0, [1, 2, 3]))   // output: 6
fmt.println(sum("", [1, 2, 3]))  // output: "123"
```

## Benchmark

## References
- [Basic Syntax](doc/syntax.md)
- [Builtin functions](doc/builtins.md)
- **Why name is Narbe?** It's from [OVERLORD](https://overlordmaruyama.fandom.com/wiki/Narberal_Gamma)
