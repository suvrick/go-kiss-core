package packets

import (
	_ "embed"
	"fmt"
	"io"
	"sort"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/schemes"
)

func NewClientPacket(w io.Writer, packetID int, payload map[string]interface{}) error {
	scheme := schemes.FindScheme(1, packetID)
	if scheme == nil {
		return fmt.Errorf("[NewClientPacket] client packet(%d) not found", packetID)
	}

	buf, err := leb128.WriteInt(packetID)
	if err != nil {
		return err
	}

	w.Write(buf) // packetID

	buf, err = leb128.WriteByte(5)
	if err != nil {
		return err
	}

	w.Write(buf) // device

	return marshal(w, scheme.Fields, payload)
}

func NewServerPacket(r io.Reader, packetID int) (map[string]interface{}, error) {

	scheme := schemes.FindScheme(0, packetID)
	if scheme == nil {
		return nil, fmt.Errorf("[NewServerPacket] server packet(%d) not found", packetID)
	}

	// if packetID == 13 {
	// 	fmt.Print("")
	// 	// TODO
	// }

	return unmarshal(r, scheme.Fields)
}

func unmarshal(r io.Reader, fields []schemes.Field) (map[string]interface{}, error) {

	var err error
	var isRequire bool
	var result = make(map[string]interface{}, 0)
	var value interface{}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Index < fields[j].Index
	})

	for _, v := range fields {

		if err != nil {
			break
		}

		isRequire = v.IsRequired

		value, err = Read(r, v.Char)

		if err != nil {
			continue
		}

		if count, ok := value.(int); ok && v.Char == 'A' {
			subSlice := make([]map[string]interface{}, 0)
			for i := 0; i < count; i++ {
				subMap, err2 := unmarshal(r, v.Children)
				if err2 != nil {
					err = err2
					break
				}

				subSlice = append(subSlice, subMap)
			}

			value = subSlice
		}

		result[v.Name] = value
	}

	if err != nil && isRequire {
		return nil, err
	}

	return result, nil
}

func Read(r io.Reader, char rune) (interface{}, error) {
	var err error
	var value interface{}
	switch char {
	case 'B':
		value, err = leb128.ReadByte(r)
	case 'I':
		value, err = leb128.ReadInt(r)
	case 'L':
		value, err = leb128.ReadLong(r)
	case 'S':
		value, err = leb128.ReadString(r)
	case 'A':
		value, err = leb128.ReadInt(r)
	default:
		err = fmt.Errorf("[Read] unsupported code %v", char)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func marshal(w io.Writer, fields []schemes.Field, payload map[string]interface{}) error {

	var err error
	var isRequire bool

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Index < fields[j].Index
	})

	for _, v := range fields {

		if err != nil {
			break
		}

		isRequire = v.IsRequired

		value, ok := payload[v.Name]
		if !ok {
			// error
			err = fmt.Errorf("[marshal] error. miss field %s", v.Name)
			continue
		}

		err = Write(w, v.Char, value)

		if childrens, ok := value.([]map[string]interface{}); ok && len(v.Children) > 0 && err == nil {
			for _, children := range childrens {
				err = marshal(w, v.Children, children)
				if err != nil {
					break
				}
			}
		}
	}

	if isRequire {
		return err
	}

	return nil
}

func Write(w io.Writer, char rune, value interface{}) error {
	var err error
	var b []byte

	switch char {
	case 'B':
		b, err = leb128.WriteByte(value)
	case 'I':
		b, err = leb128.WriteInt(value)
	case 'L':
		b, err = leb128.WriteLong(value)
	case 'S':
		b, err = leb128.WriteString(value)
	case 'A':
		if v, ok := value.([]map[string]interface{}); ok {
			b, err = leb128.WriteInt(len(v))
		}
	default:
		err = fmt.Errorf("[Write] unsupported code %v", char)
	}

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func GetByte(fieldName string, packet map[string]interface{}) (byte, bool) {
	val, ok := packet[fieldName].(byte)
	return val, ok
}

func GetInt(fieldName string, packet map[string]interface{}) (int, bool) {
	val, ok := packet[fieldName].(int)
	return val, ok
}

func GetLong(fieldName string, packet map[string]interface{}) (uint64, bool) {
	val, ok := packet[fieldName].(uint64)
	return val, ok
}

func GetString(fieldName string, packet map[string]interface{}) (string, bool) {
	val, ok := packet[fieldName].(string)
	return val, ok
}

func GetMap(fieldName string, packet map[string]interface{}) (map[string]interface{}, bool) {
	val, ok := packet[fieldName].(map[string]interface{})
	return val, ok
}

func GetMapArray(fieldName string, packet map[string]interface{}) ([]map[string]interface{}, bool) {
	val, ok := packet[fieldName].([]map[string]interface{})
	return val, ok
}
