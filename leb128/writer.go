package leb128

import (
	"reflect"
)

func Append(result []byte, value interface{}) []byte {

	if result == nil {
		result = make([]byte, 0)
	}

	switch value.(type) {
	case int8, int16, int, int32, int64:
		result = append(result, appendInt(result, toInt64(value))...)
	case uint8, uint16, uint, uint32, uint64, float32, float64:
		result = append(result, appendUint(result, toUInt64(value))...)
	case string:
		result = append(result, writeString(value.(string))...)
	}

	return result
}

func toUInt64(number interface{}) uint64 {
	v := reflect.ValueOf(number)

	if v.CanUint() {
		return v.Uint()
	}

	return 0
}

func toInt64(number interface{}) int64 {
	v := reflect.ValueOf(number)

	if v.CanInt() {
		return v.Int()
	}

	return 0
}

func writeString(str string) []byte {
	bytes := []byte(str)
	length := len(bytes)
	result := make([]byte, 0)
	result = append(result, appendInt(result, int64(length))...)
	result = append(result, bytes...)
	return result
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
