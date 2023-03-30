package packets

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/suvrick/go-kiss-core/types"
	"golang.org/x/exp/slices"
)

/*
LOGIN(4): "LBBS,BSIIBSBSBS"

	client.Login{
		ID:0x64771e0,
		NetType:0x1e,
		DeviceType:0x5,
		Key:"7b0a077a088b9e5169bcfc0bf2ee9ae8",
		OAuth:0x1,
		AccessToken:"5d5b1908c2bae78eeb199db47fc327ac935ccfbd914a38",
		Referrer:0x0,
		Tag:0x0,
		FieldInt:0x0,
		FieldString:"",
		RoomLanguage:0x0,
		FieldString2:"",
		Gender:0x0,
		Captcha:""
	}
*/
func TestLoginPacket(t *testing.T) {

	format := "LBBS,BSIIBSBSBS"
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

	buffer, err := marshal([]rune(format), data)

	fmt.Printf("%v, %v\n", buffer, err)

	want := []byte{224, 227, 157, 50, 30, 5, 32, 55, 98, 48, 97, 48, 55, 55, 97, 48, 56, 56, 98, 57, 101, 53, 49, 54, 57, 98, 99, 102, 99, 48, 98, 102, 50, 101, 101, 57, 97, 101, 56, 1, 46, 53, 100, 53, 98, 49, 57, 48, 56, 99, 50, 98, 97, 101, 55, 56, 101, 101, 98, 49, 57, 57, 100, 98, 52, 55, 102, 99, 51, 50, 55, 97, 99, 57, 51, 53, 99, 99, 102, 98, 100, 57, 49, 52, 97, 51, 56, 0, 0, 0, 0, 0, 0, 0, 0}

	fmt.Printf("Test: %v\n", slices.Compare(buffer, want) == 0)

}

func TestCustomDataMarshal(t *testing.T) {

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

	buffer, err := marshal([]rune(format), data)

	fmt.Printf("%v, %v\n", buffer, err)

	want := []byte{251, 0, 2, 7, 97, 97, 97, 97, 97, 97, 97, 7, 98, 98, 98, 98, 98, 98, 98, 3, 222, 1, 205, 2, 188, 3, 7, 99, 99, 99, 99, 99, 99, 99, 7, 122, 122, 122, 122, 122, 122, 122, 1, 171, 4}

	fmt.Printf("Test: %v\n", slices.Compare(buffer, want) == 0)
}
