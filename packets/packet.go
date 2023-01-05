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
	var char rune
	var err error
	var cPointer int
	var dPointer int
	var isArray bool

	for cPointer < len(formats) {

		if err != nil {
			if skip {
				return nil
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
			isArray = true
			cPointer++
			continue
		case ']':
			isArray = false
			cPointer++
			continue
		}

		if dPointer >= len(data) {
			err = fmt.Errorf("marshal: miss value of index: %d for data: %v", dPointer, data)
			continue
		}

		value := data[dPointer]

		if isArray {
			subFormats, err2 := getSubFormat(cPointer, formats)
			err = err2
			if err == nil {
				if newData, ok := value.([]interface{}); ok {
					if err = leb128.Write(w, len(newData)); err == nil {
						err = marshal(w, subFormats, newData)
						cPointer += len(subFormats) - 1 //прибавляем вычитанную длину группы
					}
				} else {
					err = fmt.Errorf("marshal: params %T not cast to slice", value)
				}
			}
		} else {
			switch char {
			case 'B':
				err = leb128.Write(w, value)
			case 'I':
				err = leb128.Write(w, value)
			case 'L':
				err = leb128.Write(w, value)
			case 'S':
				err = leb128.Write(w, value)
			case 'A':
				// TODO: imnplement for array ...
			default:
				err = fmt.Errorf("marshal: unsupported code %v", char)
				continue
			}
		}

		dPointer++
		cPointer++
	}

	return err
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
