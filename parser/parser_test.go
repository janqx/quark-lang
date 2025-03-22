package parser_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/janqx/quark-lang/v1/parser"
)

func TestParser(t *testing.T) {
	source, err := os.ReadFile(`C:\Users\Administrator\Desktop\projects\quark-lang\example\oop.jango`)
	if err != nil {
		panic(err)
	}
	p := parser.NewParser("<temp file>", []byte(source))
	chunk, err := p.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(chunk.String())
}
