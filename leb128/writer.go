package leb128

import (
	"fmt"

	"github.com/suvrick/go-kiss-core/until"
)

func Compress(value interface{}) ([]byte, error) {
	buffer := make([]byte, 0)
	switch t := value.(type) {
	// case nil:
	// 	buffer = appendUint(buffer, 0)
	case int8, int16, int, int32, int64:
		i64, err := until.ToInt64(value)
		if err != nil {
			return buffer, fmt.Errorf("leb128 error. {%s}", err)
		}
		buffer = appendInt(buffer, i64)
	case uint8, uint16, uint, uint32, uint64:
		ui64, err := until.ToUInt64(value)
		if err != nil {
			return buffer, fmt.Errorf("leb128 error. {%s}", err)
		}
		buffer = appendUint(buffer, ui64)
	case float32, float64:
		f64, err := until.ToFloat64(value)
		if err != nil {
			return buffer, fmt.Errorf("leb128 error. {%s}", err)
		}
		buffer = appendInt(buffer, int64(f64))
	case string:
		str, err := writeString(value)
		if err != nil {
			return buffer, fmt.Errorf("leb128 error. {%s}", err)
		}
		buffer = append(buffer, str...)
	default:
		return buffer, fmt.Errorf("leb128 error. unsupported type: %T", t)
	}
	return buffer, nil
}

func writeString(value interface{}) ([]byte, error) {
	s, err := until.ToString(value)
	if err != nil {
		return nil, err
	}
	bytes := []byte(s)
	length := len(bytes)
	result := make([]byte, 0)
	result = append(result, appendInt(result, int64(length))...)
	result = append(result, bytes...)
	//	result = append(result, 0)
	return result, nil
}

// AppendUleb128 appends v to b using unsigned LEB128 encoding.
func appendUint(b []byte, v uint64) []byte {

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
func appendInt(b []byte, v int64) []byte {
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
