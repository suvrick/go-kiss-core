package packets

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jcalabro/leb128"
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

func NewClientPacket(id ClientPacketType, data any, w io.Writer) error {

	p := GetClientScheme(id)

	err := WriteByte(w, types.B(id))
	if err != nil {
		return err
	}

	err = WriteByte(w, types.B(5))
	if err != nil {
		return err
	}

	return marshal([]rune(p.Format), data, w)
}

func marshal(formats []rune, data any, w io.Writer) error {

	var skip bool
	var char rune
	var err error
	var charPointer int
	var dataPointer int
	var subFormat []rune
	var values []any

	if v, ok := data.([]any); ok {
		values = v
	} else {
		values = []any{data}
	}

	if len(formats) > 0 && len(values) == 0 {
		return fmt.Errorf("[marshal] Empty data. Want format: %s", string(formats))
	}

	for charPointer < len(formats) {

		if err != nil {
			break
		}

		char = formats[charPointer]

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
					err = WriteInt(w, types.I(len(subData)))
					if err == nil {
						// обходим массив
						for _, v := range subData {
							err = marshal(subFormat, v, w)
							if err != nil {
								goto end
							}
						}

						charPointer += len(subFormat) + 2 // +2 []
						dataPointer++
					}
				} else {
					err = fmt.Errorf("[marshal] fail cast %T to []any", values[dataPointer])
				}
			}
			continue
		}

		if dataPointer >= len(values) {
			err = fmt.Errorf("[marshal] miss value of index: %d for data: %v", dataPointer, values)
			continue
		}

		value := values[dataPointer]

		err = writeData(char, value, w)

		dataPointer++
		charPointer++
	}
end:
	if skip {
		err = nil
	}

	return err
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

func writeData(char rune, value any, w io.Writer) error {
	switch char {
	case 'B':
		return WriteByte(w, value)
	case 'I':
		return WriteInt(w, value)
	case 'L':
		return WriteLong(w, value)
	case 'S':
		return WriteString(w, value)
	default:
		return fmt.Errorf("[writeData]: unsupported code %v", char)
	}
}

func NewServerPacket(id ServerPacketType, r io.Reader) ([]any, error) {
	p := GetServerScheme(id)
	if p == nil {
		return nil, fmt.Errorf("[NewServerPacket] not found server packet (%d)", id)
	}

	values, err := unmarshal([]rune(p.Format), r)

	if p.Name == "INFO" && len(values) > 1 {
		p2 := GetServerScheme(502)
		r2 := bytes.NewReader(values[0].(types.A))
		return unmarshal([]rune(p2.Format), r2)
	}

	return values, err
}

func unmarshal(formats []rune, r io.Reader) ([]any, error) {

	var skip bool
	var char rune
	var err error
	var charPointer int
	var subFormat []rune
	var value any
	var values []any

	for charPointer < len(formats) {

		if err != nil {
			break
		}

		char = formats[charPointer]

		switch char {
		case ',':
			skip = true
			charPointer++
			continue
		case '[':
			subFormat, err = getSubFormat(charPointer, formats)
			subValue := make([]any, 0)
			if err == nil {
				var lenArr types.I
				lenArr, err = ReadInt(r)
				for lenArr > 0 {
					value, err = unmarshal(subFormat, r)
					if err != nil {
						goto end
					}

					if len(subFormat) == 1 {
						subValue = append(subValue, value.([]any)[0])
					} else {
						subValue = append(subValue, value)
					}
					lenArr--
				}

				values = append(values, subValue)
				charPointer += len(subFormat) + 2
			}
			continue
		}

		value, err = readData(char, r)
		if err != nil {
			continue
		}

		values = append(values, value)
		charPointer++
	}

end:
	if skip {
		err = nil
	}

	return values, err
}

func readData(char rune, r io.Reader) (any, error) {
	switch char {
	case 'B':
		return ReadByte(r)
	case 'I':
		return ReadInt(r)
	case 'L':
		return ReadLong(r)
	case 'S':
		return ReadString(r)
	case 'A':
		return ReadByteArray(r)
	default:
		return nil, fmt.Errorf("[readData]: unsupported code %v", char)
	}
}

func GetClientScheme(id ClientPacketType) *Scheme {
	for _, p := range schemes {
		if p.ID == int(id) && p.Type == 1 {
			return &p
		}
	}

	return nil
}

func GetServerScheme(id ServerPacketType) *Scheme {
	for _, p := range schemes {
		if p.ID == int(id) && p.Type == 0 {
			return &p
		}
	}

	return nil
}

func ReadByteArray(r io.Reader) (types.A, error) {
	value, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("[ReadByteArray] fail cast %T of types.A", value)
	} else {
		buffer := make([]byte, value)
		_, err = r.Read(buffer)
		if err != nil {
			return nil, fmt.Errorf("[ReadByteArray] fail cast %T of types.A", value)
		}
		return types.A(buffer), nil
	}
}

func ReadByte(r io.Reader) (types.B, error) {
	value, err := leb128.DecodeU64(r)
	if err != nil {
		return 0, fmt.Errorf("[ReadByte] fail cast %T of types.B", value)
	} else {
		return types.B(value), nil
	}
}

func ReadInt(r io.Reader) (types.I, error) {
	value, err := leb128.DecodeS64(r)
	if err != nil {
		return 0, fmt.Errorf("[ReadInt] fail cast %T of types.I", value)
	} else {
		return types.I(value), nil
	}
}

func ReadLong(r io.Reader) (types.L, error) {
	value, err := leb128.DecodeU64(r)
	if err != nil {
		return 0, fmt.Errorf("[ReadLong] fail cast %T of types.L", value)
	} else {
		return types.L(value), nil
	}
}

func ReadString(r io.Reader) (types.S, error) {
	stringLen, err := leb128.DecodeU64(r)
	if err != nil {
		return types.S(""), fmt.Errorf("[ReadString] fail cast %T of types.S", stringLen)
	} else {
		str := make([]byte, stringLen)
		_, err := r.Read(str)
		if err != nil {
			return types.S(""), fmt.Errorf("[ReadString] very much len of string %v", stringLen)
		} else {
			return types.S(str), nil
		}
	}
}

func WriteByte(w io.Writer, value any) error {
	if v, ok := value.(types.B); !ok {
		return fmt.Errorf("[WriteByte] fail cast %T to types.B", value)
	} else {
		_, err := w.Write(leb128.EncodeU64(uint64(v)))
		if err != nil {
			return fmt.Errorf("[WriteByte] fail cast %T to types.B", value)
		}
	}
	return nil
}

func WriteInt(w io.Writer, value any) error {
	if v, ok := value.(types.I); !ok {
		return fmt.Errorf("[WriteInt] fail cast %T to types.I", value)
	} else {
		_, err := w.Write(leb128.EncodeS64(int64(v)))
		if err != nil {
			return fmt.Errorf("[WriteInt] fail cast %T to types.I", value)
		}
	}
	return nil
}

func WriteLong(w io.Writer, value any) error {
	if v, ok := value.(types.L); !ok {
		return fmt.Errorf("[WriteLong] fail cast %T to types.L", value)
	} else {
		_, err := w.Write(leb128.EncodeU64(uint64(v)))
		if err != nil {
			return fmt.Errorf("[WriteLong] fail cast %T to types.L", value)
		}
	}
	return nil
}

func WriteString(w io.Writer, value any) error {
	if v, ok := value.(types.S); !ok {
		return fmt.Errorf("[WriteString] fail cast %T to types.S", value)
	} else {
		_, err := w.Write(leb128.EncodeU64(uint64(len(v))))
		if err != nil {
			return fmt.Errorf("[WriteString] fail cast %T to types.S", value)
		}
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("[WriteString] fail cast %T to types.S", value)
		}
	}
	return nil
}
