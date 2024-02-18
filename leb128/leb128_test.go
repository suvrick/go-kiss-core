package leb128

import (
	"log"
	"reflect"
	"testing"
)

type s struct {
	B  byte
	F  int64
	UF uint64
	S  string
	BN string
}

func TestMarshal(t *testing.T) {

	b := []byte{32, 3}

	inst := s{
		B:  1,
		F:  2,
		UF: 3,
		S:  "4",
		BN: "5",
	}

	r := []byte{}

	r, err := WriteByte(r, inst.B)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	r, err = WriteInt64(r, inst.F)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	r, err = WriteUInt64(r, inst.UF)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	r, err = WriteString(r, inst.S)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	r, err = WriteBigNumber(r, inst.BN)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if reflect.DeepEqual(b, r) {
		log.Printf("OK")
	} else {
		log.Printf("FAIL")
	}
}
