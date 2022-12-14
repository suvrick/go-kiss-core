package packets

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
)

type Scheme struct {
	ID     int    `json:"id"`
	Type   int    `json:"type"`
	Name   string `json:"name"`
	Format string `json:"format"`
}

type ServerPacketType int
type ClientPacketType int

const (
	C_LOGIN ClientPacketType = 4
)

const (
	S_LOGIN ServerPacketType = 4
)

var schemes []Scheme

//go:embed packets.json
var f []byte

func init() {
	json.Unmarshal(f, &schemes)
}

func NewClientPacket(id ClientPacketType, data []any) ([]byte, error) {

	p := getClientScheme(id)

	w := new(bytes.Buffer)

	err := marshal(w, []rune(p.Format), data)

	return w.Bytes(), err
}

func marshal(w io.ByteWriter, formats []rune, data []any) error {

	var skip bool
	var err error
	var cPoiner int

	for cPoiner < len(formats) {

		char := formats[cPoiner]

		if cPoiner >= len(data) && !skip {
			return fmt.Errorf("")
		}

		value := data[cPoiner]

		switch char {
		case ',':
			skip = true
		case '[':
		case ']':
		case 'B':
			err = leb128.Write(w, value)
		case 'I':
			err = leb128.Write(w, value)
		case 'L':
			err = leb128.Write(w, value)
		case 'S':
			err = leb128.Write(w, value)
		case 'A':
		default:
			break
		}

	}

	return err
}

func getValue() (any, bool) {}

func NewServerPacket(id ServerPacketType, r io.ByteReader) ([]any, error) {
	p := getServerScheme(id)
	return unmarshal(p, r)
}

func unmarshal(p *Scheme, r io.ByteReader) ([]any, error) {
	return nil, nil
}

func getClientScheme(id ClientPacketType) *Scheme {
	for _, p := range schemes {
		if p.ID == int(id) && p.Type == 1 {
			return &p
		}
	}

	return nil
}

func getServerScheme(id ServerPacketType) *Scheme {
	for _, p := range schemes {
		if p.ID == int(id) && p.Type == 0 {
			return &p
		}
	}

	return nil
}
