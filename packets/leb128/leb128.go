package leb128

import (
	"errors"
	"io"
	"reflect"
)

var (
	ErrMarshalClientPacket   = errors.New("error marshal client packet")
	ErrUnmarshalServerPacket = errors.New("error unmarshal server packet")
)

func Unmarshal(reader io.Reader, s interface{}) error {

	var _err error = nil
	var _uint uint64 = 0
	var _int int64 = 0

	//fmt.Printf("%v\n", reflect.TypeOf(s))

	structure := reflect.ValueOf(s)

	numfield := reflect.ValueOf(s).Elem().NumField()

	for x := 0; x < numfield; x++ {

		field := structure.Elem().Field(x)

		switch reflect.ValueOf(s).Elem().Field(x).Kind() {
		case reflect.Bool:
			if _int, _err = ReadInt(reader, 8); _err == nil {
				if _int == 0 {
					field.SetBool(false)
				} else {
					field.SetBool(true)
				}
			}
		case reflect.Int8:
			if _int, _err = ReadInt(reader, 8); _err == nil {
				field.SetInt(_int)
			}
		case reflect.Uint8:
			if _uint, _err = ReadUint(reader, 8); _err == nil {
				field.SetUint(_uint)
			}
		case reflect.Int16:
			if _int, _err = ReadInt(reader, 16); _err == nil {
				field.SetInt(_int)
			}
		case reflect.Uint16:
			if _uint, _err = ReadUint(reader, 16); _err == nil {
				field.SetUint(_uint)
			}
		case reflect.Int:
			if _int, _err = ReadInt(reader, 32); _err == nil {
				field.SetInt(_int)
			}
		case reflect.Uint:
			if _uint, _err = ReadUint(reader, 32); _err == nil {
				field.SetUint(_uint)
			}
		case reflect.Int32:
			if _int, _err = ReadInt(reader, 32); _err == nil {
				field.SetInt(_int)
			}
		case reflect.Uint32:
			if _uint, _err = ReadUint(reader, 32); _err == nil {
				field.SetUint(_uint)
			}
		case reflect.Int64:
			if _int, _err = ReadInt(reader, 64); _err == nil {
				field.SetInt(_int)
			}
		case reflect.Uint64:
			if _uint, _err = ReadUint(reader, 64); _err == nil {
				field.SetUint(_uint)
			}
		case reflect.String:
			if _uint, _err = ReadUint(reader, 16); _err == nil {

				if _uint < 0 {
					_err = ErrUnmarshalServerPacket
					continue
				}

				str := make([]byte, _uint)
				reader.Read(str)
				field.SetString(string(str))
			}
		case reflect.Slice:
			length, _ := ReadInt(reader, 32)

			struct_type := reflect.TypeOf(field.Interface()).Elem()

			slice := reflect.MakeSlice(reflect.SliceOf(struct_type), 0, 0)

			for w := 0; w < int(length); w++ {

				item := reflect.New(struct_type)

				_err = Unmarshal(reader, item.Interface())
				if _err != nil {
					continue
				}

				slice = reflect.Append(slice, item.Elem())
			}

			field.Set(slice)
		default:
			return ErrUnmarshalServerPacket
		}
	}

	if _err != nil {
		return ErrUnmarshalServerPacket
	}

	return nil
}

func Marshal(s interface{}) ([]byte, error) {

	result := make([]byte, 0)

	var err error = nil

	structure := reflect.ValueOf(s)

	if structure.IsZero() {
		return nil, nil
	}

	numfield := reflect.ValueOf(s).Elem().NumField()

	for x := 0; x < numfield; x++ {

		if err != nil {
			return nil, err
		}

		field := structure.Elem().Field(x)

		switch reflect.ValueOf(s).Elem().Field(x).Kind() {
		case reflect.Bool:
			if i, ok := field.Interface().(bool); ok {
				if i {
					result = AppendInt(result, int64(1))
				} else {
					result = AppendInt(result, int64(0))
				}
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int8:
			if i, ok := field.Interface().(int8); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint8:
			if i, ok := field.Interface().(uint8); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int16:
			if i, ok := field.Interface().(int16); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint16:
			if i, ok := field.Interface().(uint16); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int:
			if i, ok := field.Interface().(int); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint:
			if i, ok := field.Interface().(uint); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int32:
			if i, ok := field.Interface().(int32); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint32:
			if i, ok := field.Interface().(uint32); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Int64:
			if i, ok := field.Interface().(int64); ok {
				result = AppendInt(result, int64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.Uint64:
			if i, ok := field.Interface().(uint64); ok {
				result = AppendUint(result, uint64(i))
			} else {
				err = ErrMarshalClientPacket
			}
		case reflect.String:
			if i, ok := field.Interface().(string); ok {
				str := []byte(i)
				result = AppendInt(result, int64(len(str)))
				result = append(result, str...)
			} else {
				err = ErrMarshalClientPacket
			}
		default:
			err = ErrMarshalClientPacket
		}
	}

	return result, nil
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
