package leb128

import (
	"fmt"
	"testing"
)

type A struct {
	int
}

func TestDebug(t *testing.T) {
	// value := float64(35.5)
	value := A{543}
	result, err := Compress(value)
	fmt.Printf("%v, %v", result, err)
}
