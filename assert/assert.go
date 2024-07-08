package assert

import (
	"fmt"
)

func That(condition bool, msg any) any {
	if !condition {
		panic(fmt.Sprintf("ASSERTION FAILED: %s", msg))
	}

	return nil
}

func Eq[T comparable](left, right T) any {
	if left != right {
		panic(fmt.Sprintf("ASSERTION FAILED: '%v' == '%v'", left, right))
	}

	return nil
}

func EqWithMessage[T comparable](left, right T, msg any) any {
	if left != right {
		panic(fmt.Sprintf("ASSERTION FAILED: '%v' == '%v': %s", left, right, msg))
	}

	return nil
}

func NotEq[T comparable](left, right T) any {
	if left == right {
		panic(fmt.Sprintf("ASSERTION FAILED: '%v' != '%v'", left, right))
	}

	return nil
}

func NotEqWithMessage[T comparable](left, right T, msg any) any {
	if left == right {
		panic(fmt.Sprintf("ASSERTION FAILED: '%v' != '%v': %s", left, right, msg))
	}

	return nil
}

func Unreachable(msg any) any {
	panic(fmt.Sprintf("ASSERTION FAILED: code is unreachable: %s", msg))
}
