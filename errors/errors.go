package errors

import (
	"fmt"
	"os"

	"github.com/Drumstickz64/golox/token"
)

func LogCliError(msg any, exitCode int) {
	fmt.Fprintln(os.Stderr, "golox:", msg)
	os.Exit(exitCode)
}

func LogUsageMessage() {
	fmt.Fprintln(os.Stderr, "Usage: golox [script [scan|parse|run]]")
	os.Exit(64)
}

func NewBuildtimeError(line, column int, where string, msg any) error {
	return fmt.Errorf("[line %d:%d] Error%v: %v", line, column, where, msg)
}

func NewRuntimeError(tok token.Token, msg any) error {
	return fmt.Errorf("encountered a runtime error: %v\n[on %d:%d]", msg, tok.Line, tok.Column)
}
