package packets

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/suvrick/go-kiss-core/types"
)

func Test_marshal(t *testing.T) {
	buffer := make([]byte, 0)
	format := "I[SS[I]],I"
	data := []interface{}{
		types.I(123),
		[]interface{}{
			[]interface{}{
				types.S("aaaaaaa"),
				types.S("bbbbbbb"),
				[]interface{}{
					types.I(222),
					types.I(333),
					types.I(444),
				},
			},
			[]interface{}{
				types.S("ccccccc"),
				types.S("zzzzzzz"),
				[]interface{}{
					types.I(555),
				},
			},
		},
	}

	buffer, err := marshal(buffer, []rune(format), data)

	fmt.Printf("%v, %v\n", buffer, err)
}
