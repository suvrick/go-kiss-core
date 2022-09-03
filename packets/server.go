package packets

import (
	"container/list"
	"fmt"
	"io"

	"github.com/suvrick/go-kiss-core/leb128"
	"github.com/suvrick/go-kiss-core/packets/meta"
)

func CreateServerPacket(r io.Reader) *Packet {

	p := Packet{}
	p.Len, p.Error = leb128.ReadInt(r)
	p.ID, p.Error = leb128.ReadInt(r)
	p.Type, p.Error = leb128.ReadUint16(r)

	if p.Type == 43 {
		fmt.Println(p.Name)
	}

	if p.Error != nil {
		return &p
	}

	name, format, ok := meta.Instance.GetServerMeta(p.Type)

	if !ok {
		p.Error = ErrNotFoundPacket
		return &p
	}

	p.Name = name
	p.Format = format

	if p.Type == 5 {
		p.Format = "IIISBSBBIBIBIIBBIII"
	}

	p.Params, p.Error = parse(r, []byte(format))

	return &p
}

// Где обработка ошибок ???
func parse(reader io.Reader, format []byte) ([]interface{}, error) {

	q := list.New()
	current := []interface{}{}
	skip := false
	up := func() {
		q.PushBack(current)
		current = []interface{}{}
	}

	down := func() {
		t := current
		el := q.Back()
		q.Remove(el)
		current = el.Value.([]interface{})
		current = append(current, t)
	}

	// // [II][II][II] [read] -> TIMEOUTS(173) "[BI]" PARAMS: [[[32 1655133163]]] ERROR: ""
	for i := 0; i < len(format); i++ {

		char := format[i]

		if char == ',' {
			skip = true
			continue
		}

		if char == '{' {
			up()
			continue
		}

		if char == '}' {
			down()
			continue
		}

		if char == '[' {
			l, err := leb128.ReadUint16(reader)
			if err != nil {
				return current, err
			}
			format = getGroup(format, i, l)
			up()
			continue
		}

		if char == ']' {
			down()
			continue
		}

		value, err := getValue(reader, char)
		if err != nil {
			if skip {
				return current, nil
			}
			return current, err
		}

		current = append(current, value)

	}

	return current, nil
}

func getGroup(format_array []byte, index int, count uint16) []byte {

	result := make([]byte, 0)
	summater := make([]byte, 0)
	fragment := make([]byte, 0)

	if len(format_array) < index || count < 1 {
		return nil
	}

	arr1 := format_array[:index+1]
	arr2 := format_array[index+1:]

	deep := 0
	end := 0

	for i, v := range arr2 {

		if v == '[' {
			deep++
		}

		if v == ']' {
			if deep == 0 {
				end = i
				break
			}
			deep--
		}

		fragment = append(fragment, v)
	}

	for count != 0 {
		summater = append(summater, '{')
		summater = append(summater, fragment...)
		summater = append(summater, '}')
		count--
	}

	result = append(result, arr1...)
	result = append(result, summater...)
	result = append(result, arr2[end:]...)
	return result
}

// Что с обработкой ошибок?
func getValue(reader io.Reader, r byte) (interface{}, error) {

	var value interface{}
	var err error

	switch r {
	case 'B':
		value, err = leb128.ReadUint8(reader)
	case 'I':
		value, err = leb128.ReadUint32(reader)
	case 'S':
		value, err = leb128.ReadString(reader)
	case 'A':
		_, err = leb128.ReadInt16(reader)
	default:
		return nil, ErrBadSignaturePacket
	}

	return value, err
}