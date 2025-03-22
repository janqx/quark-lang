package parser_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/janqx/quark-lang/v1/parser"
)

func TestLexer(t *testing.T) {
	source, _ := os.ReadFile(`C:\Users\Administrator\Desktop\projects\quark-lang\parser\1.txt`)
	l := parser.NewLexer("", strings.NewReader(string(source)))
	start := time.Now()
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(l.Next().String())
	fmt.Println(time.Since(start).Milliseconds())
}
