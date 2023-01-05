package packets

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"
)

func Test_marshal(t *testing.T) {
	w := new(bytes.Buffer)
	format := "I[SS[I]],I"
	data := []interface{}{
		222,
		[]interface{}{
			// "aaaaaa",
			"bbbbbb",
			[]interface{}{
				1,
				2,
			},
		},
		55,
	}

	err := marshal(w, []rune(format), data)

	fmt.Printf("%v, %v\n", w.Bytes(), err)
}
