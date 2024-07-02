package reporting

import (
	"fmt"
	"log"
	"os"
)

var HadError = false

func CliError(msg any, exitCode int) {
	fmt.Fprintln(os.Stderr, "golox: ", msg)
	os.Exit(exitCode)
}

func Error(line int, msg any) {
	report(line, "", msg)
}

func report(line int, where string, msg any) {

	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %v\n", line, where, msg)
	HadError = true
}

func ImplementationError(msg any) {
	log.Fatalln("FATAL: there was an error in the language implementation: ", msg)
}
