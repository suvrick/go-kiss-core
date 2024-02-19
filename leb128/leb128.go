package leb128

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
)

var (
	x00 = big.NewInt(0x00)
	x7F = big.NewInt(0x7F)
	x80 = big.NewInt(0x80)
)

func ReadByte(r *bytes.Reader) (byte, error) {
	return r.ReadByte()
}

func ReadInt64(r *bytes.Reader) (int64, error) {
	result, err := decodeSigned(r)
	if err != nil {
		return 0, err
	}
	return result.Int64(), nil
}

func ReadUInt64(r *bytes.Reader) (uint64, error) {
	result, err := decodeUnsigned(r)
	if err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

func ReadString(r *bytes.Reader) (string, error) {
	strLen, err := ReadUInt64(r)
	if err == nil {
		buf := make([]byte, strLen)
		r.Read(buf)
		return string(buf), err
	}
	return "", err
}

func ReadBigNumber(r *bytes.Reader) (string, error) {
	big, err := decodeUnsigned(r)
	if err != nil {
		return "", err
	}
	return big.String(), nil
}

func WriteByte(b []byte, v byte) ([]byte, error) {
	return append(b, v), nil
}

func WriteInt64(b []byte, number int64) ([]byte, error) {
	big := new(big.Int).SetInt64(number)
	result, err := encodeSigned(big)
	return append(b, result...), err
}

func WriteUInt64(b []byte, number uint64) ([]byte, error) {
	big := new(big.Int).SetUint64(number)
	result, err := encodeUnsigned(big)
	return append(b, result...), err
}

func WriteString(b []byte, s string) ([]byte, error) {
	big := new(big.Int).SetUint64(uint64(len(s)))
	result, err := encodeUnsigned(big)
	b = append(b, result...)
	return append(b, []byte(s)...), err
}

func WriteBigNumber(b []byte, s string) ([]byte, error) {
	big, ok := new(big.Int).SetString(s, 10)
	if ok {
		result, err := encodeUnsigned(big)
		return append(b, result...), err
	}
	return nil, fmt.Errorf("not number %v", s)
}

func decodeSigned(r *bytes.Reader) (*big.Int, error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	l := 0
	for _, b := range bs {
		if b < 0x80 {
			if (b & 0x40) == 0 {
				*r = *bytes.NewReader(bs)
				return decodeUnsigned(r)
			}
			break
		}
		l++
	}
	if l >= len(bs) {
		return nil, fmt.Errorf("too short")
	}
	*r = *bytes.NewReader(bs[l+1:])

	v := new(big.Int)
	for i := l; i >= 0; i-- {
		v = v.Mul(v, x80)
		v = v.Add(v, big.NewInt(int64(0x80-(bs[i]&0x7F)-1)))
	}
	v = v.Mul(v, big.NewInt(-1))
	v = v.Add(v, big.NewInt(-1))
	return v, nil
}

func encodeSigned(n *big.Int) ([]byte, error) {
	v := new(big.Int).Set(n)
	neg := v.Sign() < 0
	if neg {
		v = v.Mul(v, big.NewInt(-1))
		v = v.Add(v, big.NewInt(-1))
	}
	var bs []byte
	for {
		b := byte(v.Int64() % 0x80)
		if neg {
			b = 0x80 - b - 1
		}
		v = v.Div(v, x80)
		if (neg && v.Sign() == 0 && b&0x40 != 0) ||
			(!neg && v.Sign() == 0 && b&0x40 == 0) {
			return append(bs, b), nil
		} else {
			bs = append(bs, b|0x80)
		}
	}
}

func decodeUnsigned(r *bytes.Reader) (*big.Int, error) {
	var (
		weight = big.NewInt(1)
		value  = new(big.Int)
	)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		value = value.Add(
			value,
			new(big.Int).Mul(big.NewInt(int64(b&0x7F)), weight),
		)
		weight = weight.Mul(weight, x80)
		if b < 0x80 {
			break
		}
	}
	return value, nil
}

func encodeUnsigned(n *big.Int) ([]byte, error) {
	v := new(big.Int).Set(n)
	if v.Sign() < 0 {
		return nil, fmt.Errorf("can not leb128 encode negative values")
	}
	var bs []byte
	for {
		i := new(big.Int).And(v, x7F)
		v = v.Div(v, x80)
		if v.Cmp(x00) == 0 {
			b := i.Bytes()
			if len(b) == 0 {
				return []byte{0}, nil
			}
			return append(bs, b...), nil
		} else {
			b := new(big.Int).Or(i, x80)
			bs = append(bs, b.Bytes()...)
		}
	}
}
