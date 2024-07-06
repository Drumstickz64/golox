package errors

import (
	"fmt"
	"os"
)

func LogCliError(msg any, exitCode int) {
	fmt.Fprintln(os.Stderr, "golox:", msg)
	os.Exit(exitCode)
}

func NewError(line int, where string, msg any) error {
	return fmt.Errorf("[line %d] Error%v: %v", line, where, msg)
}
