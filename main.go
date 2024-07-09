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

const (
	SCANNING_TEST_FILEPATH = "scanning_test.lox"
	PARSING_TEST_FILEPATH  = "parsing_test.lox"
)

func main() {
	if len(os.Args) == 1 {
		RunPrompt()
	} else if len(os.Args) == 2 {
		RunFile(os.Args[1])
	} else if len(os.Args) == 3 {
		if os.Args[1] == "test" {
			if os.Args[2] == "scanning" {
				TestScanning()
			} else if os.Args[2] == "parsing" {
				// TestParsing()
			} else {
				errors.LogCliError("Can currently test: 'scanning', 'parsing'", 64)
			}
		} else {
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
			break
		}

		expression, errs := Build(line)
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
			continue PromptLoop
		}

		if err := Run(expression); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func Build(source string) ([]ast.Stmt, []error) {
	scanner := scanning.NewScanner(source)
	tokens, errs := scanner.ScanTokens()
	if len(errs) > 0 {
		return nil, errs
	}
	parser := parsing.NewParser(tokens)
	expression, errs := parser.Parse()

	if len(errs) > 0 {
		return nil, errs
	}

	return expression, errs
}

func Run(statements []ast.Stmt) error {
	interpreter := interpreting.NewInterpreter()
	return interpreter.Interpret(statements)
}

func TestScanning() {
	source := LoadSource(SCANNING_TEST_FILEPATH)

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

// func TestParsing() {
// 	source := LoadSource(PARSING_TEST_FILEPATH)

// 	scanner := scanning.NewScanner(source)
// 	tokens, errs := scanner.ScanTokens()
// 	for _, err := range errs {
// 		fmt.Fprintln(os.Stderr, err)
// 	}

// 	if len(errs) > 0 {
// 		os.Exit(65)
// 	}

// 	parser := parsing.NewParser(tokens)
// 	expression, err := parser.Parse()
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(65)
// 	}

// 	fmt.Println(ast.NewPrinter().Print(expression))
// }
