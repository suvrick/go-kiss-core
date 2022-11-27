package leb128

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
)

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var (
	ErrMarshalClientPacket   = errors.New("error marshal client packet")
	ErrUnmarshalServerPacket = errors.New("error unmarshal server packet")
)

// AppendUleb128 appends v to b using unsigned LEB128 encoding.
func AppendUint(b []byte, v uint64) []byte {
	for {
		//берём 7 бит
		// 	13 -> 0 0 0 0  1 1 0 1
		// 127 -> 0 1 1 1  1 1 1 1
		c := uint8(v & 0x7f)
		//сдвигаем на 7 бит
		v >>= 7
		if v != 0 {
			// 1 0 0 0 0 0 0 0
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

// AppendSleb128 appends v to b using signed LEB128 encoding.
func AppendInt(b []byte, v int64) []byte {
	for {
		c := uint8(v & 0x7f) // берем первых 7 бит
		s := uint8(v & 0x40)
		v >>= 7 // сдвигайем на 7 бит вправо
		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			// если вошли сюда то
			c |= 0x80 // дописываем 8 бит
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

func ReadUint(r io.Reader, n uint) (uint64, error) {
	if n > 64 {
		return 0, errors.New("leb128: invalid uint")
	}
	p := make([]byte, 1)
	var res uint64
	var shift uint

	for {
		_, err := io.ReadFull(r, p)
		if err != nil {
			return 0, err
		}
		b := uint64(p[0])

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

func ReadInt(r io.Reader, n uint) (int64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
	}
	p := make([]byte, 1)
	var res int64
	var shift uint
	for {
		_, err := io.ReadFull(r, p)
		if err != nil {
			return 0, err
		}
		b := int64(p[0])
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
