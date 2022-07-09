package leb128

import (
	"errors"
)

func Write(value interface{}) []byte {

	b := make([]byte, 0)

	// default return for nil argument
	if value == nil {
		b = append(b, 0)
		return b
	}

	switch v := value.(type) {
	case int8:
		b = AppendInt(b, int64(value.(int8)))
	case uint8:
		b = AppendUint(b, uint64(value.(uint8)))
	case int16:
		b = AppendInt(b, int64(value.(int16)))
	case uint16:
		b = AppendUint(b, uint64(value.(uint16)))
	case int:
		b = AppendInt(b, int64(value.(int)))
	case uint:
		b = AppendUint(b, uint64(value.(uint)))
	case int32:
		b = AppendInt(b, int64(value.(int32)))
	case uint32:
		b = AppendUint(b, uint64(value.(uint32)))
	case int64:
		b = AppendInt(b, int64(value.(int64)))
	case uint64:
		b = AppendUint(b, uint64(value.(uint64)))
	case float32:
		b = AppendInt(b, int64(value.(float32)))
	case float64:
		b = AppendUint(b, uint64(value.(float64)))
	case string:
		str := value.(string)
		str_len := int64(len(str))
		b = AppendInt(b, str_len)

		if str_len != 0 {
			b = append(b, []byte(str)...)
		}

	default:
		_ = v
		return b
	}

	return b
}

func Append(b []byte, value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case int8:
		b = AppendInt(b, int64(value.(int8)))
	case uint8:
		b = AppendUint(b, uint64(value.(uint8)))
	case int16:
		b = AppendInt(b, int64(value.(int16)))
	case uint16:
		b = AppendUint(b, uint64(value.(uint16)))
	case int:
		b = AppendInt(b, int64(value.(int)))
	case uint:
		b = AppendUint(b, uint64(value.(uint)))
	case int32:
		b = AppendInt(b, int64(value.(int32)))
	case uint32:
		b = AppendUint(b, uint64(value.(uint32)))
	case int64:
		b = AppendInt(b, int64(value.(int64)))
	case uint64:
		b = AppendUint(b, uint64(value.(uint64)))
	case float32:
		b = AppendInt(b, int64(value.(float32)))
	case float64:
		b = AppendUint(b, uint64(value.(float64)))
	case string:
		str := value.(string)
		b = AppendInt(b, int64(len(str)))
		if len(str) != 0 {
			b = append(b, []byte(str)...)
		}
	default:
		_ = v
		return b, errors.New("bad signature")
	}

	return b, nil
}

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
