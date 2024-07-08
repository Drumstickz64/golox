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
	LogCliError("Usage: golox [script] | golox test <operation>", 64)
}

func NewBuildtimeError(line int, where string, msg any) error {
	return fmt.Errorf("[line %d] Error%v: %v", line, where, msg)
}

func NewRuntimeError(tok token.Token, msg any) error {
	return fmt.Errorf("%v\n[line  %d]", msg, tok.Line)
}
