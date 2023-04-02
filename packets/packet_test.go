package packets

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"

	"github.com/suvrick/go-kiss-core/types"
	"golang.org/x/exp/slices"
)

func TestWriteLoginPacket(t *testing.T) {
	format := "LBBS,BSIIBSBSBS"
	//            [56 0 5 4 147 129 234 214 214 233 184 128 186 1 5 1 32 101 49 100 101 55 100 54 98 49 98 57 97 49 56 101 49 50 52 51 51 49 100 49 97 56 101 55 97 54 55 48 57 1 0 0 0 0 0 0 0 0 0]
	want := []byte{224, 227, 157, 50, 30, 5, 32, 55, 98, 48, 97, 48, 55, 55, 97, 48, 56, 56, 98, 57, 101, 53, 49, 54, 57, 98, 99, 102, 99, 48, 98, 102, 50, 101, 101, 57, 97, 101, 56, 1, 46, 53, 100, 53, 98, 49, 57, 48, 56, 99, 50, 98, 97, 101, 55, 56, 101, 101, 98, 49, 57, 57, 100, 98, 52, 55, 102, 99, 51, 50, 55, 97, 99, 57, 51, 53, 99, 99, 102, 98, 100, 57, 49, 52, 97, 51, 56, 0, 0, 0, 0, 0, 0, 0, 0}
	data := []interface{}{
		types.L(105345504),
		types.B(30),
		types.B(5),
		types.S("7b0a077a088b9e5169bcfc0bf2ee9ae8"),
		types.B(1),
		types.S("5d5b1908c2bae78eeb199db47fc327ac935ccfbd914a38"),
		types.I(0),
		types.I(0),
		types.B(0),
		types.S(""),
		types.B(0),
		types.S(""),
		types.B(0),
		types.S(""),
	}

	w := bytes.NewBuffer(nil)

	err := marshal([]rune(format), data, w)

	fmt.Printf("%v, %v\n", w.Bytes(), err)

	fmt.Printf("Test: %v\n", slices.Compare(w.Bytes(), want) == 0)
}

func TestWriteCustomDataMarshal(t *testing.T) {
	format := "I[SS[I]],I"
	want := []byte{251, 0, 2, 7, 97, 97, 97, 97, 97, 97, 97, 7, 98, 98, 98, 98, 98, 98, 98, 3, 222, 1, 205, 2, 188, 3, 7, 99, 99, 99, 99, 99, 99, 99, 7, 122, 122, 122, 122, 122, 122, 122, 1, 171, 4}
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

	w := bytes.NewBuffer(nil)

	err := marshal([]rune(format), data, w)

	fmt.Printf("%v, %v\n", w.Bytes(), err)

	fmt.Printf("Test: %v\n", slices.Compare(w.Bytes(), want) == 0)
}

func TestReadCustomDataUnmarshal(t *testing.T) {
	format := "I[SS[I]],I"
	data := []byte{251, 0, 2, 7, 97, 97, 97, 97, 97, 97, 97, 7, 98, 98, 98, 98, 98, 98, 98, 3, 222, 1, 205, 2, 188, 3, 7, 99, 99, 99, 99, 99, 99, 99, 7, 122, 122, 122, 122, 122, 122, 122, 1, 171, 4, 0}
	want := []any{
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
		types.I(0),
	}

	r := bytes.NewReader(data)

	stuff, err := unmarshal([]rune(format), r)

	fmt.Printf("%v, %v\n", stuff, err)

	fmt.Printf("Stuff: %s\n", fmt.Sprintf("%v", stuff))

	fmt.Printf("Want: %s\n", fmt.Sprintf("%v", want))

	fmt.Printf("Test: %v\n", fmt.Sprintf("%v", stuff) == fmt.Sprintf("%v", want))
}
