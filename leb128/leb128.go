package leb128

import (
	"errors"
	"io"
)

const (
	maxWidth = 10
)

var (
	ErrReadByte   = errors.New("[ReadByte] fail read byte")
	ErrReadInt    = errors.New("[ReadInt] fail read int")
	ErrReadLong   = errors.New("[ReadLong] fail read long")
	ErrReadString = errors.New("[ReadString] fail read string")

	ErrWriteByte   = errors.New("[WriteByte] fail write byte")
	ErrWriteInt    = errors.New("[WriteInt] fail write int")
	ErrWriteLong   = errors.New("[WriteLong] fail write long")
	ErrWriteString = errors.New("[WriteString] fail write string")

	ErrToUint64 = errors.New("[ToUint64] upsupported type")

	ErrOverflow = errors.New("LEB128 integer overflow (was more than 8 bytes)")
)

func ReadByte(r io.Reader) (byte, error) {
	value, err := decodeU64(r, 1)
	if err != nil {
		return 0, ErrReadByte
	}
	return byte(value), nil
}

func ReadInt(r io.Reader) (int, error) {
	value, err := decodeU64(r, 4)
	if err != nil {
		return 0, ErrReadInt
	}
	return int(value), nil
}

func ReadLong(r io.Reader) (uint64, error) {
	value, err := decodeU64(r, 8)
	if err != nil {
		return 0, ErrReadLong
	}
	return value, nil
}

func ReadString(r io.Reader) (string, error) {
	if strLen, err := ReadInt(r); err == nil {
		strBuf := make([]byte, strLen)
		if _, err = r.Read(strBuf); err == nil {
			return string(strBuf), nil
		}
	}
	return "", ErrReadString
}

func WriteByte(value interface{}) ([]byte, error) {
	v, err := ToUint64(value)
	if err != nil {
		return nil, ErrWriteByte
	}
	return encodeU64(v), nil
}

func WriteInt(value interface{}) ([]byte, error) {
	v, err := ToUint64(value)
	if err != nil {
		return nil, ErrWriteInt
	}
	return encodeU64(v), nil
}

func WriteLong(value interface{}) ([]byte, error) {
	v, err := ToUint64(value)
	if err != nil {
		return nil, ErrWriteLong
	}
	return encodeU64(v), nil
}

func WriteString(value interface{}) ([]byte, error) {
	if str, ok := value.(string); ok {
		if lenStr, err := WriteInt(len(str)); err == nil {
			return append(lenStr, []byte(str)...), nil
		}
	}
	return nil, ErrWriteString
}

func ToUint64(value any) (uint64, error) {
	switch v := value.(type) {
	case int8:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case int:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case uint64:
		return v, nil
	case float32:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	}
	return 0, ErrToUint64
}

func decodeU64(r io.Reader, size int) (uint64, error) {
	var res uint64

	bit := int8(0)
	buf := make([]byte, 1)
	for i := 0; ; i++ {
		if i > maxWidth {
			return 0, ErrOverflow
		}

		if i == size {
			return uint64(uint8(buf[0])), nil
		}

		_, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		b := buf[0]

		res |= uint64(b&0x7f) << (7 * bit)

		signBit := b & 0x80
		if signBit == 0 {
			break
		}

		bit++
	}

	return res, nil
}

func decodeS64(r io.Reader, size int) (int64, error) {
	var res int64

	shift := 0
	buf := make([]byte, 1)
	for i := 0; ; i++ {
		if i > maxWidth {
			return 0, ErrOverflow
		}

		if i == size {
			return int64(int8(buf[0])), nil
		}

		_, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		b := buf[0]

		res |= int64(b&0x7f) << shift
		shift += 7

		if b&0x80 == 0 {
			if b&0x40 != 0 {
				// signed
				res |= ^0 << shift
			}
			break
		}
	}

	return res, nil
}

func encodeU64(num uint64) []byte {
	buf := make([]byte, 0)

	done := false
	for !done {
		b := byte(num & 0x7F)

		num = num >> 7
		if num == 0 {
			done = true
		} else {
			b |= 0x80
		}

		buf = append(buf, b)
	}

	return buf
}

func encodeS64(num int64) []byte {
	buf := make([]byte, 0)

	done := false
	for !done {
		//
		// From https://go.dev/ref/spec#Arithmetic_operators:
		//
		// "The shift operators implement arithmetic shifts
		// if the left operand is a signed integer and
		// logical shifts if it is an unsigned integer"
		//

		b := byte(num & 0x7F)
		num >>= 7 // arithmetic shift
		signBit := b & 0x40
		if (num == 0 && signBit == 0) ||
			(num == -1 && signBit != 0) {
			done = true
		} else {
			b |= 0x80
		}

		buf = append(buf, b)
	}

	return buf
}
