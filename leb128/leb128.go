package leb128

import (
	"errors"
	"io"
)

const (
	maxWidth = 10
)

var (
	ErrOverflow = errors.New("LEB128 integer overflow (was more than 8 bytes)")
)

func ReadSByte(r io.Reader) (int64, error) {
	return decodeS64(r, 1)
}

func ReadUByte(r io.Reader) (uint64, error) {
	return decodeU64(r, 1)
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

// DecodeU64 converts a uleb128 byte stream to a uint64. Be careful
// to ensure that your data can fit in 8 bytes.
func DecodeU64(r io.Reader) (uint64, error) {
	var res uint64

	bit := int8(0)
	buf := make([]byte, 1)
	for i := 0; ; i++ {
		if i > maxWidth {
			return 0, ErrOverflow
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

// DecodeS64 converts a sleb128 byte stream to a int64. Be careful
// to ensure that your data can fit in 8 bytes.
func DecodeS64(r io.Reader) (int64, error) {
	var res int64

	shift := 0
	buf := make([]byte, 1)
	for i := 0; ; i++ {
		if i > maxWidth {
			return 0, ErrOverflow
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

// EncodeU64 converts num to a uleb128 encoded array of bytes
func EncodeU64(num uint64) []byte {
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

// EncodeS64 converts num to a sleb128 encoded array of bytes
func EncodeS64(num int64) []byte {
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
