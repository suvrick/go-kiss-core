package leb128

import (
	"errors"
	"fmt"
	"io"
)

func ReadInt8(r io.Reader) (int8, error) {
	value, err := readInt(r, 8)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read int8")
	}
	return int8(value), nil
}

func ReadInt16(r io.Reader) (int16, error) {
	value, err := readInt(r, 16)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read int16")
	}
	return int16(value), nil
}

func ReadInt(r io.Reader) (int, error) {
	value, err := readInt(r, 32)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read int")
	}
	return int(value), nil
}

func ReadInt32(r io.Reader) (int32, error) {
	value, err := readInt(r, 32)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read int32")
	}
	return int32(value), nil
}

func ReadInt64(r io.Reader) (int64, error) {
	value, err := readInt(r, 64)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read int64")
	}
	return value, nil
}

func ReadUint8(r io.Reader) (uint8, error) {
	value, err := readUint(r, 8)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read uin8")
	}
	return uint8(value), nil
}

func ReadUint16(r io.Reader) (uint16, error) {
	value, err := readUint(r, 16)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read uin16")
	}
	return uint16(value), nil
}

func ReadUint(r io.Reader) (uint, error) {
	value, err := readUint(r, 32)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read uint")
	}
	return uint(value), nil
}

func ReadUint32(r io.Reader) (uint32, error) {
	value, err := readUint(r, 32)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read uint32")
	}
	return uint32(value), nil
}

func ReadUint64(r io.Reader) (uint64, error) {
	value, err := readUint(r, 64)
	if err != nil {
		return 0, fmt.Errorf("leb128 error. Can`t read uint64")
	}
	return value, nil
}

func ReadString(r io.Reader) (string, error) {
	length, err := ReadUint16(r)
	if err != nil {
		return "", fmt.Errorf("leb128 error. Can`t read string")
	}
	data := make([]byte, length)
	_, err = r.Read(data)
	if err != nil {
		return "", fmt.Errorf("leb128 error. Can`t read string")
	}
	return string(data), nil
}

func readUint(r io.Reader, n uint) (uint64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
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

func readInt(r io.Reader, n uint) (int64, error) {
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
