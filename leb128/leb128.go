package leb128

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/suvrick/go-kiss-core/types"
)

var (
	ErrMarshalClientPacket   = errors.New("error marshal client packet")
	ErrUnmarshalServerPacket = errors.New("error unmarshal server packet")
)

func Skip(v interface{}) bool {
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		if value, ok := t.Field(i).Tag.Lookup("pack"); ok {
			if value == "skip" {
				return true
			}
		}
	}
	return false
}

func Unmarshal(reader io.Reader, s interface{}) (err error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			err = fmt.Errorf("[unmarshal] %v", r)
		}
	}()

	value := reflect.ValueOf(s)

	err = unmarshal(reader, value)

	return
}

func unmarshal(reader io.Reader, value reflect.Value) error {

	var _err error = nil
	var _uint uint64 = 0
	var _int int64 = 0

	switch value.Kind() {
	case reflect.Bool:
		if _int, _err = ReadInt(reader, 8); _err == nil {
			if _int == 0 {
				value.SetBool(false)
			} else {
				value.SetBool(true)
			}
		}
	case reflect.Int8:
		if _int, _err = ReadInt(reader, 8); _err == nil {
			value.SetInt(_int)
		}
	case reflect.Uint8:
		if _uint, _err = ReadUint(reader, 8); _err == nil {
			value.SetUint(_uint)
		}
	case reflect.Int16:
		if _int, _err = ReadInt(reader, 16); _err == nil {
			value.SetInt(_int)
		}
	case reflect.Uint16:
		if _uint, _err = ReadUint(reader, 16); _err == nil {
			value.SetUint(_uint)
		}
	case reflect.Int:
		if _int, _err = ReadInt(reader, 32); _err == nil {
			value.SetInt(_int)
		}
	case reflect.Uint:
		if _uint, _err = ReadUint(reader, 32); _err == nil {
			value.SetUint(_uint)
		}
	case reflect.Int32:
		if _int, _err = ReadInt(reader, 32); _err == nil {
			value.SetInt(_int)
		}
	case reflect.Uint32:
		if _uint, _err = ReadUint(reader, 32); _err == nil {
			value.SetUint(_uint)
		}
	case reflect.Int64:
		if _int, _err = ReadInt(reader, 64); _err == nil {
			value.SetInt(_int)
		}
	case reflect.Uint64:
		if _uint, _err = ReadUint(reader, 64); _err == nil {
			value.SetUint(_uint)
		}
	case reflect.String:
		_uint, _err = ReadUint(reader, 16)
		if _err == nil {
			str := make([]byte, _uint)
			reader.Read(str)
			value.SetString(string(str))
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			unmarshal(reader, getField(value, i))
		}
	case reflect.Slice:
		length, _ := ReadUint(reader, 16)
		struct_type := reflect.TypeOf(value.Interface()).Elem()
		slice := reflect.MakeSlice(reflect.SliceOf(struct_type), 0, 0)
		for w := 0; w < int(length); w++ {
			item := reflect.New(struct_type)
			_err = unmarshal(reader, item)
			if _err != nil {
				continue
			}
			slice = reflect.Append(slice, item.Elem())
		}
		value.Set(slice)
	case reflect.Pointer:
		if value.Pointer() != 0 {
			unmarshal(reader, value.Elem())
		}
	default:
	}

	return nil
}

func getField(v reflect.Value, i int) reflect.Value {
	val := v.Field(i)
	if val.Kind() == reflect.Interface && !val.IsNil() {
		val = val.Elem()
	}
	return val
}

func Marshal(s interface{}) (result []byte, err error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			err = fmt.Errorf("[marshal] %v", r)
		}
	}()

	value := reflect.ValueOf(s)

	return marshal(value)
}

func marshal(v reflect.Value) ([]byte, error) {

	result := make([]byte, 0)
	var err error

	switch v.Kind() {
	case reflect.Bool:
		if i, ok := v.Interface().(types.B); ok {
			result = AppendInt(result, int64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Int8:
		if i, ok := v.Interface().(types.B); ok {
			result = AppendInt(result, int64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Uint8:
		if i, ok := v.Interface().(types.B); ok {
			result = AppendUint(result, uint64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Int16:
		if i, ok := v.Interface().(types.B); ok {
			result = AppendInt(result, int64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Uint16:
		if i, ok := v.Interface().(types.B); ok {
			result = AppendUint(result, uint64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Int:
		if i, ok := v.Interface().(types.I); ok {
			result = AppendInt(result, int64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Uint:
		if i, ok := v.Interface().(types.I); ok {
			result = AppendUint(result, uint64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Int32:
		if i, ok := v.Interface().(types.I); ok {
			result = AppendInt(result, int64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Uint32:
		if i, ok := v.Interface().(types.I); ok {
			result = AppendUint(result, uint64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Int64:
		if i, ok := v.Interface().(types.I); ok {
			result = AppendInt(result, int64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Uint64:
		if i, ok := v.Interface().(types.I); ok {
			result = AppendUint(result, uint64(i))
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.String:
		if i, ok := v.Interface().(types.S); ok {
			str := []byte(i)
			result = AppendInt(result, int64(len(str)))
			result = append(result, str...)
		} else {
			err = ErrMarshalClientPacket
		}
	case reflect.Slice:

		lenn := v.Len()

		result = AppendInt(result, int64(lenn))

		for i := 0; i < lenn; i++ {

			val := v.Index(i)

			res, _err := marshal(val)
			if err != nil {
				return nil, _err
			}
			result = append(result, res...)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			r, _err := marshal(v.Field(i))
			if err != nil {
				err = _err
			} else {
				result = append(result, r...)
			}
		}
	case reflect.Pointer:
		if v.Pointer() != 0 {
			r, _err := marshal(v.Elem())
			if err != nil {
				err = _err
			} else {
				result = append(result, r...)
			}
		}
	default:
		err = ErrMarshalClientPacket
	}

	return result, err
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
