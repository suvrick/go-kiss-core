package packets

import (
	_ "embed"
	"encoding/json"
	"fmt"

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

func NewClientPacket(id ClientPacketType, data any) ([]byte, error) {

	p := getClientScheme(id)

	return marshal([]rune(p.Format), data)
}

func marshal(formats []rune, data any) ([]byte, error) {

	var skip bool
	var char rune
	var err error
	var charPointer int
	var dataPointer int
	var subFormat []rune
	var values []any
	var buffer []byte

	if v, ok := data.([]any); ok {
		values = v
	} else {
		values = []any{data}
	}

	if len(formats) > 0 && len(values) == 0 {
		return nil, fmt.Errorf("[marshal] Empty data. Want format: %s", string(formats))
	}

	for charPointer < len(formats) {

		if err != nil {
			break
		}

		char = formats[charPointer]
		r := string(char)
		_ = r

		switch char {
		case ',':
			skip = true
			charPointer++
			continue
		case '[':
			// "I[SS[I]]" -> "I SS[I] SS[I]" -> "I SS III SS II"
			subFormat, err = getSubFormat(charPointer, formats)
			if err == nil {
				// конвертим к []any
				if subData, ok := values[dataPointer].([]any); ok {
					// Записываем длину массива
					buffer, err = leb128.WriteInt(buffer, types.I(len(subData)))
					if err == nil {
						// обходим массив
						for _, v := range subData {
							newBuffer := make([]byte, 0)
							newBuffer, err = marshal(subFormat, v)
							if err != nil {
								return buffer, err
							}
							buffer = append(buffer, newBuffer...)
						}
					}
				}
			}

			charPointer += len(subFormat) + 2
			dataPointer++
			continue
		}

		if dataPointer >= len(values) {
			err = fmt.Errorf("[marshal] miss value of index: %d for data: %v", dataPointer, values)
			continue
		}

		value := values[dataPointer]

		buffer, err = writeData(char, value, buffer)

		dataPointer++
		charPointer++
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
	//case 'A':
	// TODO: imnplement for array ...
	default:
		err = fmt.Errorf("[writeData]: unsupported code %v", char)
	}

	return buffer, err
}

func NewServerPacket(id ServerPacketType, buffer []byte) ([]any, error) {
	p := getServerScheme(id)
	return unmarshal([]rune(p.Format), buffer)
}

func unmarshal(format []rune, buffer []byte) ([]any, error) {

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
