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
		types.B(123),
	}

	buffer, err := marshal(buffer, []rune(format), data)

	fmt.Printf("%v, %v\n", buffer, err)
}
