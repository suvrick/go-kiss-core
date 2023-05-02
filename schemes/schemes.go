package schemes

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

var schms []Scheme

// //go:embed packets.json
// var f []byte

// func init() {
// 	json.Unmarshal(f, &schemes)
// }

func init() {
	s := []Scheme{
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
					Name:       "frame_type",
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
		{
			PacketID:     5,
			PacketType:   0,
			PacketName:   "INFO",
			PacketFormat: "ILBSBIIISBSSBBIIIIBBIIII",
			Fields: []Field{
				{
					Index:      0,
					Name:       "length",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "players",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "game_id",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      1,
							Name:       "login_id",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      2,
							Name:       "frame_type",
							Char:       'B',
							IsRequired: true,
						},
						{
							Index:      3,
							Name:       "name",
							Char:       'S',
							IsRequired: true,
						},
					},
				},
			},
		},
	}

	SetSchemes(s)
}

func SetSchemes(s []Scheme) {
	schms = s
}

func ClearSchemes() {
	schms = nil
}

func FindScheme(packetType int, packetID int) *Scheme {
	for _, v := range schms {
		if packetType == v.PacketType && packetID == v.PacketID {
			return &v
		}
	}
	return nil
}
