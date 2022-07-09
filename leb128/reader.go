package leb128

import (
	"errors"
	"io"
)

func ReadUint(r io.Reader, n uint) (uint64, error) {
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
