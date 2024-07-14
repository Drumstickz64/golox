package main

import (
	"bufio"
	goerrors "errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Drumstickz64/golox/ast"
	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/interpreting"
	"github.com/Drumstickz64/golox/parsing"
	"github.com/Drumstickz64/golox/scanning"
)

func main() {
	if len(os.Args) == 1 {
		RunPrompt()
	} else if len(os.Args) == 2 {
		RunFile(os.Args[1])
	} else if len(os.Args) == 3 {
		switch os.Args[2] {
		case "scan":
			TestScanning(os.Args[1])
		case "parse":
			TestParsing(os.Args[1])
		case "run":
			RunFile(os.Args[1])
		default:
			errors.LogUsageMessage()
		}
	} else {
		errors.LogUsageMessage()
	}
}

func RunFile(path string) {
	source := LoadSource(path)

	expression, errs := Build(source)
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}

	if len(errs) > 0 {
		os.Exit(65)
	}

	if err := Run(expression); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(70)
	}
}

func LoadSource(path string) string {
	source, err := os.ReadFile(path)
	if err != nil {
		errors.LogCliError(err, 66)
	}

	return string(source)
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

PromptLoop:
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if goerrors.Is(err, io.EOF) {
				break
			}

			errors.LogCliError(err, 74)
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue PromptLoop
		}

		statements, errs := Build(line)
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
			continue PromptLoop
		}

		if err := Run(statements); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func Build(source string) ([]ast.Stmt, []errors.BuildError) {
	scanner := scanning.NewScanner(source)
	tokens, errs := scanner.ScanTokens()
	if len(errs) > 0 {
		return nil, errs
	}
	parser := parsing.NewParser(tokens)
	statements, errs := parser.Parse()

	if len(errs) > 0 {
		return nil, errs
	}

	return statements, errs
}

func Run(statements []ast.Stmt) error {
	interpreter := interpreting.NewInterpreter()
	return interpreter.Interpret(statements)
}

func TestScanning(pth string) {
	source := LoadSource(pth)

	scanner := scanning.NewScanner(source)
	tokens, errs := scanner.ScanTokens()

	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}

	if len(errs) > 0 {
		os.Exit(65)
	}

	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

func TestParsing(pth string) {
	panic("Parsing test currently unavailable")
	// source := LoadSource(pth)

	// scanner := scanning.NewScanner(source)
	// tokens, errs := scanner.ScanTokens()
	// for _, err := range errs {
	// 	fmt.Fprintln(os.Stderr, err)
	// }

	// if len(errs) > 0 {
	// 	os.Exit(65)
	// }

	// parser := parsing.NewParser(tokens)
	// expression, err := parser.Parse()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(65)
	// }

	// fmt.Println(ast.NewPrinter().Print(expression))
}
