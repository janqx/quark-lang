## 定义变量与基础类型展示
```javascript
a = true    // boolean
b = 567       // integer
c = -7834.12  // float
d = "hello，世界！"   // string
e = [a, b, c, d] // list
f = {name: "tom", age: 20, misc: e} // dict
g = fn(a, b) { // function
    return a + b
}
```

## if
```javascript
if cond1 {

} else if cond2 {

} else if cond3 {

} else {

}
```

## for
```javascript
for {
    // 无限循环，相当于C语言中的while(1)
}

for cond {
    // 单条件循环，相当于C语言中的while
}

for i := 0; i < 10; i += 1 {
    // 迭代循环，相当于C语言中的for
}
```

## function
```javascript
fn add(a, b) {
    c = a + b
    return c
}
```

## import and export
```javascript
// a.ng
math := import("math")
fn calc(a, b) {
    c = math.pow(a + b, 2)
    return c
}
export {
    add: add
}

// b.ng
a := import("a") // amount to `a := import("a.ng")`
print(a.add(5, 3)) // output: 64
```

