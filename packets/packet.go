package packets

import (
	_ "embed"
	"fmt"
	"io"
	"sort"

	"github.com/suvrick/go-kiss-core/leb128"
)

type Field struct {
	ID         int
	Name       string
	Index      int
	Char       rune
	IsRequired bool
	Parent     *Field
	Children   []Field
}

type Scheme struct {
	PacketID     int     `json:"id"`
	PacketType   int     `json:"type"`
	PacketName   string  `json:"name"`
	PacketFormat string  `json:"format"`
	Fields       []Field `json:"fields"`
}

var schemes []Scheme

// //go:embed packets.json
// var f []byte

// func init() {
// 	json.Unmarshal(f, &schemes)
// }

func init() {
	schemes = []Scheme{
		{
			PacketID:     4,
			PacketType:   1,
			PacketName:   "LOGIN",
			PacketFormat: "LBBS,BSIIBBS",
			Fields: []Field{
				{
					Index:      0,
					Name:       "login_id",
					Char:       'L',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "net_type",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      3,
					Name:       "auth_key",
					Char:       'S',
					IsRequired: true,
				},
				{
					Index:      2,
					Name:       "device",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     4,
			PacketType:   0,
			PacketName:   "LOGIN",
			PacketFormat: "B,IIII[B]IIIISBBIIBBS",
			Fields: []Field{
				{
					Index:      0,
					Name:       "status",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "game_id",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      2,
					Name:       "balance",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      3,
					Name:       "inviter_id",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      4,
					Name:       "logout_time",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      5,
					Name:       "flags",
					Char:       'A',
					IsRequired: false,
					Children: []Field{
						{
							Index:      0,
							Name:       "flag",
							Char:       'B',
							IsRequired: false,
						},
					},
				},
				{
					Index:      6,
					Name:       "games_count",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      7,
					Name:       "kisses_daily_count",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      8,
					Name:       "last_payment_time",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      9,
					Name:       "subscribe_expires",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      10,
					Name:       "params",
					Char:       'S',
					IsRequired: false,
				},
				{
					Index:      11,
					Name:       "sex_set",
					Char:       'B',
					IsRequired: false,
				},
				{
					Index:      12,
					Name:       "tutorial",
					Char:       'B',
					IsRequired: false,
				},
				{
					Index:      13,
					Name:       "tag",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      14,
					Name:       "server_time",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      15,
					Name:       "first_login",
					Char:       'B',
					IsRequired: false,
				},
				{
					Index:      16,
					Name:       "is_top_player",
					Char:       'B',
					IsRequired: false,
				},
				{
					Index:      16,
					Name:       "photos_hash",
					Char:       'S',
					IsRequired: false,
				},
			},
		},
	}
}

func FindScheme(packetType int, packetID int) *Scheme {
	for _, v := range schemes {
		if packetType == v.PacketType && packetID == v.PacketID {
			return &v
		}
	}
	return nil
}

func NewClientPacket(w io.Writer, packetID int, payload map[string]interface{}) error {
	scheme := FindScheme(1, packetID)
	if scheme == nil {
		return fmt.Errorf("[NewClientPacket] client packet(%d) not found", packetID)
	}

	buf, err := leb128.WriteInt(packetID)
	if err != nil {
		return err
	}

	w.Write(buf) // packetID

	buf, err = leb128.WriteByte(5)
	if err != nil {
		return err
	}

	w.Write(buf) // device

	return marshal(w, scheme.Fields, payload)
}

func NewServerPacket(r io.Reader, packetID int) (map[string]interface{}, error) {
	scheme := FindScheme(0, packetID)
	if scheme == nil {
		return nil, fmt.Errorf("[NewServerPacket] server packet(%d) not found", packetID)
	}

	return unmarshal(r, scheme.Fields)
}

func unmarshal(r io.Reader, fields []Field) (map[string]interface{}, error) {

	var err error
	var isRequire bool
	var result = make(map[string]interface{}, 0)
	var value interface{}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Index < fields[j].Index
	})

	for _, v := range fields {

		if err != nil {
			break
		}

		isRequire = v.IsRequired

		value, err = Read(r, v.Char)

		if err != nil {
			continue
		}

		if count, ok := value.(int); ok && v.Char == 'A' {
			subSlice := make([]interface{}, 0)
			for i := 0; i < count; i++ {
				subMap, err2 := unmarshal(r, v.Children)
				if err2 != nil {
					err = err2
					break
				}

				subSlice = append(subSlice, subMap)
			}

			value = subSlice
		}

		result[v.Name] = value
	}

	if isRequire {
		return nil, err
	}

	return result, nil
}

func Read(r io.Reader, char rune) (interface{}, error) {
	var err error
	var value interface{}
	switch char {
	case 'B':
		value, err = leb128.ReadByte(r)
	case 'I':
		value, err = leb128.ReadInt(r)
	case 'L':
		value, err = leb128.ReadLong(r)
	case 'S':
		value, err = leb128.ReadString(r)
	case 'A':
		value, err = leb128.ReadInt(r)
	default:
		err = fmt.Errorf("[Read] unsupported code %v", char)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func marshal(w io.Writer, fields []Field, payload map[string]interface{}) error {

	var err error
	var isRequire bool

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Index < fields[j].Index
	})

	for _, v := range fields {

		if err != nil {
			break
		}

		isRequire = v.IsRequired

		value, ok := payload[v.Name]
		if !ok {
			// error
			err = fmt.Errorf("[marshal] error. miss field %s", v.Name)
			continue
		}

		err = Write(w, v.Char, value)

		if childrens, ok := value.([]map[string]interface{}); ok && len(v.Children) > 0 && err == nil {
			for _, children := range childrens {
				err = marshal(w, v.Children, children)
				if err != nil {
					break
				}
			}
		}
	}

	if isRequire {
		return err
	}

	return nil
}

func Write(w io.Writer, char rune, value interface{}) error {
	var err error
	var b []byte

	switch char {
	case 'B':
		b, err = leb128.WriteByte(value)
	case 'I':
		b, err = leb128.WriteInt(value)
	case 'L':
		b, err = leb128.WriteLong(value)
	case 'S':
		b, err = leb128.WriteString(value)
	case 'A':
		if v, ok := value.([]map[string]interface{}); ok {
			b, err = leb128.WriteInt(len(v))
		}
	default:
		err = fmt.Errorf("[Write] unsupported code %v", char)
	}

	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}
