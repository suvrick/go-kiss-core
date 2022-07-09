package packets

import (
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
	p.Buffer, p.Error = leb128.Append(p.Buffer, p_type) // packet type
	p.Buffer, p.Error = leb128.Append(p.Buffer, 4)      // device type

	if p.Error != nil {
		return p
	}

	if p.Type == 8 {
		p.Format = "[I]I,I"
	}

	arr := load([]byte(p.Format), params)
	p.Buffer = append(p.Buffer, arr...)

	return p
}

func (p *Packet) GetBuffer(msgID int64) []byte {

	a := leb128.Write(msgID)
	b := len(p.Buffer) + len(a)
	c := leb128.Write(b)

	data := make([]byte, 0)
	data = append(data, c...)        // итоговая длина пакета
	data = append(data, a...)        // ID сообщения
	data = append(data, p.Buffer...) // данные
	return data

}

func load(format []byte, params []interface{}) []byte {

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
			value := setValue(1, 'B')
			current = append(current, value...)
			continue
		}

		if char == ']' {
			continue
		}

		value := setValue(next(), char)
		current = append(current, value...)
	}

	return current
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

func setValue(v interface{}, code byte) []byte {

	switch code {
	case 'B', 'I', 'L', 'F', 'S':
		return leb128.Write(v)
	}

	return nil
}
