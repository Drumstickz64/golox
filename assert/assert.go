package assert

import (
	"fmt"
)

func That(condition bool, msg string) {
	if !condition {
		panic(fmt.Sprintln("ASSERTION FAILED:", msg))
	}
}

func Eq[T comparable](left, right T) {
	if left != right {
		panic(fmt.Sprintf("ASSERTION FAILED: '%v' == '%v'", left, right))
	}
}

func EqWithMessage[T comparable](left, right T, msg string) {
	if left != right {
		panic(fmt.Sprintf("ASSERTION FAILED: '%v' == '%v': %s", left, right, msg))
	}
}

func NotEq[T comparable](left, right T) {
	if left == right {
		panic(fmt.Sprintf("ASSERTION FAILED: '%v' != '%v'", left, right))
	}
}

func NotEqWithMessage[T comparable](left, right T, msg string) {
	if left == right {
		panic(fmt.Sprintf("ASSERTION FAILED: '%v' != '%v': %s", left, right, msg))
	}
}
