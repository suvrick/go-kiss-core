package packets

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"

	"github.com/suvrick/go-kiss-core/packets/leb128"
	"github.com/suvrick/go-kiss-core/types"
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

func NewClientPacket(id ClientPacketType, values any) (data []byte, err error) {

	p := getClientScheme(id)

	return marshal(data, []rune(p.Format), values)
}

func marshal(buffer []byte, formats []rune, values1 any) ([]byte, error) {

	var skip bool
	var char rune
	var err error
	var cPointer int
	var dPointer int
	var subFormat []rune
	var values []any

	if v, ok := values1.([]any); ok {
		values = v
	} else {
		values = []any{values1}
	}

	for cPointer < len(formats) {

		if err != nil {
			break
		}

		char = formats[cPointer]
		r := string(char)
		_ = r

		switch char {
		case ',':
			skip = true
			cPointer++
			continue
		case '[':
			// "I[SS[I]]" -> "I SS[I] SS[I]" -> "I SS III SS II"
			subFormat, err = getSubFormat(cPointer, formats)
			if err == nil {
				// конвертим к []any
				if data, ok := values[dPointer].([]any); ok {
					// Записываем длину массива
					buffer, err = leb128.WriteInt(buffer, types.I(len(data)))
					if err == nil {
						// обходим массив
						for _, v := range data {
							buffer, err = marshal(buffer, subFormat, v)
							if err != nil {
								return buffer, err
							}
						}
					}
				}
			}

			cPointer += len(subFormat) + 2
			dPointer++
			continue
		}

		if dPointer >= len(values) {
			err = fmt.Errorf("marshal: miss value of index: %d for data: %v", dPointer, values)
			continue
		}

		value := values[dPointer]

		buffer, err = writeData(char, value, buffer)

		dPointer++
		cPointer++
	}

	if skip {
		err = nil
	}

	return buffer, err
}

func getSubFormat(index int, formats []rune) ([]rune, error) {
	r := make([]rune, 0)
	dep := 0
	for i := index + 1; i < len(formats); i++ {
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

	if dep != -1 {
		return nil, fmt.Errorf("[getSubFormat] invalid packet format")
	}

	return r, nil
}

func writeData(char rune, value any, buffer []byte) ([]byte, error) {

	var err error

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
		err = fmt.Errorf("[writeData]: unsupported code %v", char)
	}

	return buffer, err
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
