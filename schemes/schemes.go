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

		// Client

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
			PacketID:     8,
			PacketType:   1,
			PacketName:   "REQUEST",
			PacketFormat: "[I]I,I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "players",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "player_id",
							Char:       'I',
							IsRequired: true,
						},
					},
				},
				{
					Index:      1,
					Name:       "mask1",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      2,
					Name:       "mask2",
					Char:       'I',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     11,
			PacketType:   1,
			PacketName:   "GAME_REWARDS_GET",
			PacketFormat: "I,I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "reward_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "count",
					Char:       'I',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     21,
			PacketType:   1,
			PacketName:   "MOVE",
			PacketFormat: "IB",
			Fields: []Field{
				{
					Index:      0,
					Name:       "target_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "destination",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     26,
			PacketType:   1,
			PacketName:   "BOTTLE_PLAY",
			PacketFormat: "B,B",
			Fields: []Field{
				{
					Index:      0,
					Name:       "type",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "language",
					Char:       'B',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     27,
			PacketType:   1,
			PacketName:   "BOTTLE_LEAVE",
			PacketFormat: "",
		},
		{
			PacketID:     28,
			PacketType:   1,
			PacketName:   "BOTTLE_ROLL",
			PacketFormat: "I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "speed",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     29,
			PacketType:   1,
			PacketName:   "BOTTLE_KISS",
			PacketFormat: "B",
			Fields: []Field{
				{
					Index:      0,
					Name:       "answer",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     30,
			PacketType:   1,
			PacketName:   "BOTTLE_SAVE",
			PacketFormat: "I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "target_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     31,
			PacketType:   1,
			PacketName:   "BOTTLE_KICK",
			PacketFormat: "I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "player_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     32,
			PacketType:   1,
			PacketName:   "PASSION_PASS_RESET",
			PacketFormat: "",
		},
		{
			PacketID:     33,
			PacketType:   1,
			PacketName:   "BOTTLE_WAITER_JOIN",
			PacketFormat: "",
		},
		{
			PacketID:     34,
			PacketType:   1,
			PacketName:   "BOTTLE_WAITER_LEAVE",
			PacketFormat: "",
		},
		{
			PacketID:     61,
			PacketType:   1,
			PacketName:   "RECEIVE_DAILY_BONUS",
			PacketFormat: "",
		},
		{
			PacketID:     202,
			PacketType:   1,
			PacketName:   "BOTTLE_MOVE",
			PacketFormat: "I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "room_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},

		// Server

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
		{
			PacketID:     7,
			PacketType:   0,
			PacketName:   "BALANCE",
			PacketFormat: "I,B",
			Fields: []Field{
				{
					Index:      0,
					Name:       "balance",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "reason",
					Char:       'B',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     10,
			PacketType:   0,
			PacketName:   "CONTEST_ITEMS",
			PacketFormat: "B,B",
			Fields: []Field{
				{
					Index:      0,
					Name:       "status",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "item_type",
					Char:       'B',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     13,
			PacketType:   0,
			PacketName:   "REWARDS",
			PacketFormat: "[II]",
			Fields: []Field{
				{
					Index:      0,
					Name:       "rewards",
					Char:       'A',
					IsRequired: false,
					Children: []Field{
						{
							Index:      0,
							Name:       "id",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      1,
							Name:       "count",
							Char:       'I',
							IsRequired: true,
						},
					},
				},
			},
		},
		{
			PacketID:     17,
			PacketType:   0,
			PacketName:   "BONUS",
			PacketFormat: "B",
			Fields: []Field{
				{
					Index:      0,
					Name:       "daily_bonus_type",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     20,
			PacketType:   0,
			PacketName:   "VIP",
			PacketFormat: "I,B",
			Fields: []Field{
				{
					Index:      0,
					Name:       "time",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "reason",
					Char:       'B',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     22,
			PacketType:   0,
			PacketName:   "MOVE",
			PacketFormat: "I,B",
			Fields: []Field{
				{
					Index:      0,
					Name:       "status",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "room_id",
					Char:       'B',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     24,
			PacketType:   0,
			PacketName:   "BOTTLE_PLAY_DENIED",
			PacketFormat: "B",
			Fields: []Field{
				{
					Index:      0,
					Name:       "status",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     25,
			PacketType:   0,
			PacketName:   "BOTTLE_ROOM",
			PacketFormat: "III[I][I]",
			Fields: []Field{
				{
					Index:      0,
					Name:       "room_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "table_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      2,
					Name:       "bottle_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      3,
					Name:       "players",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "player_id",
							Char:       'I',
							IsRequired: true,
						},
					},
				},
				{
					Index:      4,
					Name:       "spectators",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "spectator_id",
							Char:       'I',
							IsRequired: true,
						},
					},
				},
			},
		},
		{
			PacketID:     26,
			PacketType:   0,
			PacketName:   "BOTTLE_JOIN",
			PacketFormat: "IB",
			Fields: []Field{
				{
					Index:      0,
					Name:       "inner_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "index",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     27,
			PacketType:   0,
			PacketName:   "BOTTLE_LEAVE",
			PacketFormat: "I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "inner_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     28,
			PacketType:   0,
			PacketName:   "BOTTLE_LEADER",
			PacketFormat: "I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "leader_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     29,
			PacketType:   0,
			PacketName:   "BOTTLE_ROLL",
			PacketFormat: "II,II",
			Fields: []Field{
				{
					Index:      0,
					Name:       "leader_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "rolled_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      2,
					Name:       "speed",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      3,
					Name:       "time",
					Char:       'I',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     30,
			PacketType:   0,
			PacketName:   "BOTTLE_KISS",
			PacketFormat: "IB",
			Fields: []Field{
				{
					Index:      0,
					Name:       "player_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "answer",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     31,
			PacketType:   0,
			PacketName:   "BOTTLE_TABLE",
			PacketFormat: "II",
			Fields: []Field{
				{
					Index:      0,
					Name:       "actor_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "table_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     32,
			PacketType:   0,
			PacketName:   "BOTTLE_BOTTLE",
			PacketFormat: "II",
			Fields: []Field{
				{
					Index:      0,
					Name:       "actor_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "bottle_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     33,
			PacketType:   0,
			PacketName:   "BOMB_EXPLODE",
			PacketFormat: "IIBB",
			Fields: []Field{
				{
					Index:      0,
					Name:       "source_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "catcher_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      2,
					Name:       "is_new",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      3,
					Name:       "room_type",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     34,
			PacketType:   0,
			PacketName:   "BOMB_EXPLODE",
			PacketFormat: "IB",
			Fields: []Field{
				{
					Index:      0,
					Name:       "player_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "room_type",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     35,
			PacketType:   0,
			PacketName:   "BOTTLE_ENTER",
			PacketFormat: "",
		},
		{
			PacketID:     37,
			PacketType:   0,
			PacketName:   "CHAT_MESSAGE",
			PacketFormat: "BIIS,II",
			Fields: []Field{
				{
					Index:      0,
					Name:       "chat_type",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "message_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      2,
					Name:       "player_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      3,
					Name:       "message",
					Char:       'S',
					IsRequired: true,
				},
				{
					Index:      4,
					Name:       "complain_target",
					Char:       'I',
					IsRequired: false,
				},
				{
					Index:      5,
					Name:       "complain_data",
					Char:       'I',
					IsRequired: false,
				},
			},
		},
		{
			PacketID:     38,
			PacketType:   0,
			PacketName:   "CHAT_WHISPER",
			PacketFormat: "ISBIB",
			Fields: []Field{
				{
					Index:      0,
					Name:       "player_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "message",
					Char:       'S',
					IsRequired: true,
				},
				{
					Index:      2,
					Name:       "direction",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      3,
					Name:       "history_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      4,
					Name:       "is_new",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     206,
			PacketType:   0,
			PacketName:   "REWARD_STATE",
			PacketFormat: "III",
			Fields: []Field{
				{
					Index:      0,
					Name:       "hearts",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "gifts",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      2,
					Name:       "time",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     212,
			PacketType:   0,
			PacketName:   "PLAYER_ROOM_TYPE",
			PacketFormat: "IB",
			Fields: []Field{
				{
					Index:      0,
					Name:       "player_id",
					Char:       'B',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "room_type",
					Char:       'B',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     218,
			PacketType:   0,
			PacketName:   "PLAYERS_KISSES",
			PacketFormat: "[II]",
			Fields: []Field{
				{
					Index:      0,
					Name:       "players_kisses",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "player_id",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      1,
							Name:       "kisses",
							Char:       'I',
							IsRequired: true,
						},
					},
				},
			},
		},
		{
			PacketID:     230,
			PacketType:   0,
			PacketName:   "PLAYERS_VIEW",
			PacketFormat: "[IIB]",
			Fields: []Field{
				{
					Index:      0,
					Name:       "viewers",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "player_id",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      1,
							Name:       "time",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      2,
							Name:       "is_unocked",
							Char:       'B',
							IsRequired: true,
						},
					},
				},
			},
		},
		{
			PacketID:     232,
			PacketType:   0,
			PacketName:   "UPDATE_INFO",
			PacketFormat: "I",
			Fields: []Field{
				{
					Index:      0,
					Name:       "player_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     304,
			PacketType:   0,
			PacketName:   "CAPTCHA",
			PacketFormat: "",
		},
		{
			PacketID:     307,
			PacketType:   0,
			PacketName:   "KISS_PRIORITY",
			PacketFormat: "[II]",
			Fields: []Field{
				{
					Index:      0,
					Name:       "kiss_priorities",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "player_id",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      1,
							Name:       "count",
							Char:       'I',
							IsRequired: true,
						},
					},
				},
			},
		},
		{
			PacketID:     308,
			PacketType:   0,
			PacketName:   "KICKS_KICKS",
			PacketFormat: "[III]",
			Fields: []Field{
				{
					Index:      0,
					Name:       "kicks",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "target_id",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      1,
							Name:       "actor_id",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      2,
							Name:       "kick_left",
							Char:       'I',
							IsRequired: true,
						},
					},
				},
			},
		},
		{
			PacketID:     309,
			PacketType:   0,
			PacketName:   "KICK_SAVE",
			PacketFormat: "II",
			Fields: []Field{
				{
					Index:      0,
					Name:       "target_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "actor_id",
					Char:       'I',
					IsRequired: true,
				},
			},
		},
		{
			PacketID:     310,
			PacketType:   0,
			PacketName:   "BALANCE_ITEMS",
			PacketFormat: "[BII]",
			Fields: []Field{
				{
					Index:      0,
					Name:       "balance_items",
					Char:       'A',
					IsRequired: true,
					Children: []Field{
						{
							Index:      0,
							Name:       "type",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      1,
							Name:       "common",
							Char:       'I',
							IsRequired: true,
						},
						{
							Index:      2,
							Name:       "daily",
							Char:       'I',
							IsRequired: true,
						},
					},
				},
			},
		},
		{
			PacketID:     315,
			PacketType:   0,
			PacketName:   "REWARD_GOT",
			PacketFormat: "II",
			Fields: []Field{
				{
					Index:      0,
					Name:       "owner_id",
					Char:       'I',
					IsRequired: true,
				},
				{
					Index:      1,
					Name:       "reward_id",
					Char:       'I',
					IsRequired: true,
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
