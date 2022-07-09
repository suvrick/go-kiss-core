package packets

import (
	"fmt"

	"github.com/suvrick/go-kiss-core/leb128"
)

var c_packets map[uint64]Packet

func SetClientPakets(p *map[uint64]Packet) {
	c_packets = *p
}

func GetClientPacket(id uint64) (Packet, bool) {
	p, ok := c_packets[id]
	return p, ok
}

func CreateClientPacket(p_type uint64, params ...interface{}) Packet {

	p, ok := GetClientPacket(p_type)
	if !ok {
		p = Packet{
			Type:   p_type,
			Params: params,
			Error:  ErrNotFoundPacket,
		}
		return p
	}

	p.Params = params

	p.Buffer = make([]byte, 0)
	data, err := leb128.Compress(p_type) // packet type
	if err != nil {
		p.Error = err
		return p
	}

	p.Buffer = append(p.Buffer, data...)

	data, err = leb128.Compress(4) // device type
	if err != nil {
		p.Error = err
		return p
	}

	p.Buffer = append(p.Buffer, data...)

	// FIX ME
	if p.Type == 8 {
		p.Format = "[I]I,I"
	}

	data, err = load([]byte(p.Format), params)
	if err != nil {
		p.Error = err
		return p
	}

	p.Buffer = append(p.Buffer, data...)
	return p
}

func (p *Packet) GetBuffer(msgID int64) ([]byte, error) {

	a, err := leb128.Compress(msgID)
	if err != nil {
		return nil, err
	}

	b := len(p.Buffer) + len(a)

	c, err := leb128.Compress(b)
	if err != nil {
		return nil, err
	}

	data := make([]byte, 0)
	data = append(data, c...)        // итоговая длина пакета
	data = append(data, a...)        // ID сообщения
	data = append(data, p.Buffer...) // данные
	return data, nil

}

func load(format []byte, params []interface{}) ([]byte, error) {

	next := nextParam(params)
	current := []byte{}

	for i := 0; i < len(format); i++ {
		char := format[i]

		if char == ',' {
			continue
		}

		if char == '{' {
			continue
		}

		if char == '}' {
			continue
		}

		if char == '[' {
			value, err := setValue(1, 'B')
			if err != nil {
				return current, err
			}

			current = append(current, value...)
			continue
		}

		if char == ']' {
			continue
		}

		value, err := setValue(next(), char)
		if err != nil {
			return current, err
		}

		current = append(current, value...)
	}

	return current, nil
}

func nextParam(params []interface{}) func() interface{} {
	i := -1
	return func() interface{} {
		i++

		if i > len(params)-1 {
			return nil
		}

		return params[i]
	}
}

func setValue(v interface{}, code byte) ([]byte, error) {
	switch code {
	case 'B', 'I', 'L', 'F', 'S':
		return leb128.Compress(v)
	default:
		return nil, fmt.Errorf("unsupported code %v", code)
	}
}
