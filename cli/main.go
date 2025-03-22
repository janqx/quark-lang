package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/janqx/quark-lang/v1"
	"github.com/janqx/quark-lang/v1/stdlib"
)

const (
	REPL_PROMPT = ">> "
)

var (
	flagShowVersion bool
	flagShowHelp    bool
	flagCmd         string
)

func repl() {
	ctx := quark.NewContext(quark.ModeREPL, stdlib.LoadModules())
	script := quark.NewScript(ctx)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(REPL_PROMPT)
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "failed to scan input")
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.ToLower(line) == ".exit" {
			break
		}
		result, err := script.RunString(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if result != nil && result != quark.Null {
			fmt.Println(quark.ToString(result))
		}
	}
}

func run(filename string) {
	ctx := quark.NewContext(quark.ModeNormal, stdlib.LoadModules())
	script := quark.NewScript(ctx)
	if err := script.RunFile(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func execute(source string) {
	ctx := quark.NewContext(quark.ModeNormal, stdlib.LoadModules())
	script := quark.NewScript(ctx)
	if _, err := script.RunString(source); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func main() {
	flag.BoolVar(&flagShowVersion, "version", false, "show version information")
	flag.BoolVar(&flagShowHelp, "help", false, "show help information")
	flag.StringVar(&flagCmd, "c", "", "execute string")
	flag.Parse()

	if flagShowHelp {
		_, executable := filepath.Split(os.Args[0])
		fmt.Printf("Usage: %s [file] [options]\nOptions:\n", executable)
		flag.PrintDefaults()
		os.Exit(0)
	} else if flagShowVersion {
		fmt.Printf("Quark v%d.%d.%d\nrepository: %s", quark.VersionMajor, quark.VersionMinor, quark.VersionPatch, "github.com/janqx/quark-lang")
		os.Exit(0)
	}

	if flagCmd != "" {
		execute(flagCmd)
		os.Exit(0)
	}

	filename := flag.Arg(0)
	if filename == "" {
		repl()
	} else {
		run(filename)
	}
}
