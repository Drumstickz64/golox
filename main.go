package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Drumstickz64/golox/reporting"
	"github.com/Drumstickz64/golox/scanning"
)

func main() {
	if len(os.Args) == 1 {
		RunPrompt()
	} else if len(os.Args) == 2 {
		RunFile(os.Args[1])
	} else {
		reporting.CliError("Usage: golox [script]", 64)
	}
}

func RunFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		reporting.CliError(err, 1)
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

			reporting.CliError(err, 1)
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
