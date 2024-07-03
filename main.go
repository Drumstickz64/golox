package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Drumstickz64/golox/ast"
	"github.com/Drumstickz64/golox/reporting"
	"github.com/Drumstickz64/golox/scanning"
	"github.com/Drumstickz64/golox/token"
)

func main() {
	expression := ast.Binary{
		Left: &ast.Grouping{
			Expression: &ast.Binary{
				Left: &ast.Literal{Value: 1},
				Operator: token.Token{
					Kind:   token.PLUS,
					Lexeme: "+",
				},
				Right: &ast.Literal{Value: 2},
			},
		},
		Operator: token.Token{
			Kind:   token.STAR,
			Lexeme: "*",
		},
		Right: &ast.Grouping{
			Expression: &ast.Binary{
				Left: &ast.Literal{Value: 4},
				Operator: token.Token{
					Kind:   token.MINUS,
					Lexeme: "-",
				},
				Right: &ast.Literal{Value: 3},
			},
		},
	}

	fmt.Println(ast.RPNPrinter{}.Print(&expression))

	// if len(os.Args) == 1 {
	// 	RunPrompt()
	// } else if len(os.Args) == 2 {
	// 	RunFile(os.Args[1])
	// } else {
	// 	reporting.CliError("Usage: golox [script]", 64)
	// }
}

func RunFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		reporting.CliError(err, 66)
	}

	Run(string(content))

	if reporting.HadError {
		os.Exit(65)
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			reporting.CliError(err, 74)
		}

		line = strings.TrimSpace(line)

		if line == "" {
			break
		}

		Run(line)

		reporting.HadError = false
	}
}

func Run(source string) {
	scanner := scanning.NewScanner(source)
	tokens := scanner.ScanTokens()

	if !reporting.HadError {
		for _, token := range tokens {
			fmt.Printf("%v\n", token)
		}
	}
}
