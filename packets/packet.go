package packets

import (
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
	S_LOGIN ServerPacketType = 4
)

var schemes []Scheme

//go:embed packets.json
var f []byte

func init() {
	json.Unmarshal(f, &schemes)
}

func NewClientPacket(id ClientPacketType, values []any) (data []byte, err error) {

	p := getClientScheme(id)

	return marshal(data, []rune(p.Format), values)
}

func marshal(buffer []byte, formats []rune, values []any) ([]byte, error) {

	var skip bool
	var char rune
	var err error
	var cPointer int
	var dPointer int
	//var isArray bool

	for cPointer < len(formats) {

		if err != nil {
			if skip {
				return nil, nil
			} else {
				break
			}
		}

		char = formats[cPointer]

		switch char {
		case ',':
			skip = true
			cPointer++
			continue
		case '[':
			//isArray = true
			cPointer++
			continue
		case ']':
			//isArray = false
			cPointer++
			continue
		}

		if dPointer >= len(values) {
			err = fmt.Errorf("marshal: miss value of index: %d for data: %v", dPointer, values)
			continue
		}

		value := values[dPointer]

		switch char {
		case 'B':
			buffer, err = leb128.WriteByte(buffer, value)
		case 'I':
			buffer, err = leb128.WriteInt(buffer, value)
		case 'L':
			buffer, err = leb128.WriteLong(buffer, value)
		case 'S':
			buffer, err = leb128.WriteString(buffer, value)
		case 'A':
			// TODO: imnplement for array ...
		default:
			err = fmt.Errorf("marshal: unsupported code %v", char)
			continue
		}
	}

	dPointer++
	cPointer++

	return buffer, err
}

func getSubFormat(index int, formats []rune) ([]rune, error) {
	r := make([]rune, 0)
	dep := 0
	for i := index; i < len(formats); i++ {
		c := formats[i]

		if c == '[' {
			dep++
		}

		if c == ']' {
			dep--
			if dep == -1 {
				break
			}
		}

		r = append(r, c)
	}

	if dep != 0 {
		return nil, fmt.Errorf("[getSubFormat] invalid packet format")
	}

	return r, nil
}

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
