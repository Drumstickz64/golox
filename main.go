package main

import (
	"bufio"
	goerrors "errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Drumstickz64/golox/ast/expr"
	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/parsing"
	"github.com/Drumstickz64/golox/scanning"
)

func main() {
	if len(os.Args) == 1 {
		RunPrompt()
	} else if len(os.Args) == 2 {
		RunFile(os.Args[1])
	} else {
		errors.LogCliError("Usage: golox [script]", 64)
	}
}

func RunFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		errors.LogCliError(err, 66)
	}

	if errs := Run(string(content)); errs != nil {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(65)
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

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

		if err := Run(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func Run(source string) []error {
	// ========== scanning test ==========

	// scanner := scanning.NewScanner(source)
	// tokens, errs := scanner.ScanTokens()

	// if len(errs) > 0 {
	// 	return errs
	// }

	// for _, tok := range tokens {
	// 	fmt.Println(tok)
	// }

	// ========== parsing test ==========

	scanner := scanning.NewScanner(source)
	tokens, errs := scanner.ScanTokens()
	parser := parsing.NewParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errs
	}

	fmt.Println(expr.NewPrinter().Print(expression))

	return nil
}
