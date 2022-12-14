package leb128

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// type Number interface {
// 	bool | int8 | *int8 | int16 | *int16 | int32 | *int32 | int64 | *int64
// 	uint8 | *uint8 | uint16 | *uint16 | uint32 | *uint32 | uint64 | *uint64
// 	float32 | *float32 | float64 | *float64
// }

func Write(buffer io.ByteWriter, v any) (err error) {

	var value any

	switch x := v.(type) {
	case int8:
		value = int64(x)
	case *int8:
		value = int64(*x)
	case uint8:
		value = uint64(x)
	case *uint8:
		value = uint64(*x)
	case int16:
		value = int64(x)
	case *int16:
		value = int64(*x)
	case uint16:
		value = uint64(x)
	case *uint16:
		value = uint64(*x)
	case int:
		value = int64(x)
	case *int:
		value = int64(*x)
	case uint:
		value = uint64(x)
	case *uint:
		value = uint64(*x)
	case int32:
		value = int64(x)
	case *int32:
		value = int64(*x)
	case uint32:
		value = uint64(x)
	case *uint32:
		value = uint64(*x)
	case int64:
		value = int64(x)
	case *uint64:
		value = uint64(*x)
	case uint64:
		value = uint64(x)
	case *int64:
		value = int64(*x)
	case float32:
		value = int64(x)
	case *float32:
		value = int64(*x)
	case float64:
		value = int64(x)
	case *float64:
		value = int64(*x)
	case bool:
		if x {
			value = uint64(1)
		} else {
			value = int64(0)
		}
	case *bool:
		if *x {
			value = uint64(1)
		} else {
			value = int64(0)
		}
	case string:
		value = int64(len(x))
	case *string:
		value = int64(len(*x))
	default:
		err = fmt.Errorf("leb128: unsupported type %T", x)
		return
	}

	switch x := value.(type) {
	case int64:
		WriteInt(buffer, x)
	case uint64:
		WriteUint(buffer, x)
	}

	if s, ok := v.(string); ok {
		if value.(int64) > 0 {
			for _, c := range []byte(s) {
				err = buffer.WriteByte(c)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func WriteInt(w io.ByteWriter, v int64) {
	for {
		c := uint8(v & 0x7f) // берем первых 7 бит
		s := uint8(v & 0x40)
		v >>= 7 // сдвигайем на 7 бит вправо
		if (v != int64(-1) || s == 0) && (v != 0 || s != 0) {
			c |= 0x80 // дописываем 8 бит
		}
		w.WriteByte(c)
		if c&0x80 == 0 {
			break
		}
	}
}

func WriteUint(w io.ByteWriter, v uint64) {
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		w.WriteByte(c)
		if c&0x80 == 0 {
			break
		}
	}
}

func ReadString(r io.ByteReader) (result string, err error) {

	var slen int64
	slen, err = ReadInt(r, 32)
	if err != nil {
		return
	}

	str := strings.Builder{}
	var b byte

	for slen > 0 {
		b, err = r.ReadByte()
		if err != nil {
			return
		}
		str.WriteByte(b)
		slen--
	}

	result = str.String()
	return
}

func ReadBool(r io.ByteReader) (result bool, err error) {
	v, e := ReadUint(r, 8)
	if e != nil {
		err = e
		return
	}

	if v != 0 {
		result = true
	}

	return
}

func ReadUint(r io.ByteReader, n uint) (uint64, error) {
	if n > 64 {
		return 0, errors.New("leb128: invalid uint")
	}

	var res uint64
	var shift uint

	for {
		p, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		b := uint64(p)

		if n == 8 {
			return b, nil
		}

		switch {
		case b < 1<<7 && b < 1<<n:
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		default:
			return 0, errors.New("leb128: invalid uint")
		}
	}
}

func ReadInt(r io.ByteReader, n uint) (int64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
	}
	var res int64
	var shift uint
	for {
		p, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		b := int64(p)
		switch {
		case b < 1<<6 && uint64(b) < uint64(1<<(n-1)):
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<6 && b < 1<<7 && uint64(b)+1<<(n-1) >= 1<<7:
			res += (1 << shift) * (b - 1<<7)
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		default:
			return 0, errors.New("leb128: invalid int")
		}
	}
}
