package packets

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"
)

func Test_marshal(t *testing.T) {
	w := new(bytes.Buffer)
	format := "[SS],I"
	data := []interface{}{
		[]interface{}{
			"aaaaaa",
			"bbbbbb",
		},
	}

	err := marshal(w, []rune(format), data)

	fmt.Printf("%v, %v\n", w.Bytes(), err)
}
