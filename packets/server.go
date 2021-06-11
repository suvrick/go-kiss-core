package packets

import (
	"io"
	"strings"

	"github.com/go-interpreter/wagon/wasm/leb128"
)

type ServerPacketType uint16

const (
	//"status:B, inner_id:I, balance:I"
	LOGIN_SERVER ServerPacketType = 4
	//"skip:I, skip:I, skip:I, net:I, type:B, name:S, sex:B, tag:B, referrer:I, bday:I, photo:S, photo_byte:B, profile:S"
	INFO_SERVER ServerPacketType = 5
	//"bottles:I, reason:B"
	BALANCE_SERVER ServerPacketType = 7
	//"skip:I, id:I, count:I, json:S"
	GAME_REWARDS_SERVER ServerPacketType = 13
	//"can_collect:B, day:B"
	BONUS_SERVER ServerPacketType = 62
	// [skip:I, rewardId:I]
	ROULETTE ServerPacketType = 262
)

var s_format = map[ServerPacketType]string{
	LOGIN_SERVER:        "status:B, inner_id:I, balance:I",
	INFO_SERVER:         "skip:I, skip:I, skip:I, net:L, type:B, name:S, sex:B, tag:B, referrer:I, bday:I, photo:S, photo_byte:B, profile:S",
	BALANCE_SERVER:      "bottles:I, reason:B",
	GAME_REWARDS_SERVER: "skip:I, id:I, count:I, json:S",
	BONUS_SERVER:        "can_collect:B, day:B",
	ROULETTE:            "skip:I, rewardId:I",
}

func getServerFormat(t ServerPacketType) string {

	line, ok := s_format[t]

	if !ok {
		return ""
	}

	return line
}

func CreateServerPacket(t ServerPacketType, format string, r io.Reader) map[string]interface{} {

	var _ ServerPacketType = t

	format = strings.ReplaceAll(format, " ", "")
	words := strings.Split(format, ",")

	if words == nil {
		return nil
	}

	result := make(map[string]interface{}, 0)
	for _, word := range words {

		word = strings.ReplaceAll(word, "[", "")
		word = strings.ReplaceAll(word, "]", "")

		s := strings.Split(word, ":")

		if len(s) != 2 {
			continue
		}

		name := s[0]
		code := s[1]

		if name == "skip" {
			leb128.ReadVarUint32(r)
			continue
		}

		switch code {
		case "B", "I":
			value, _ := leb128.ReadVarUint32(r)
			result[name] = value
			break
		case "L":
			value, _ := leb128.ReadVarint64(r)
			result[name] = value
			break
		case "S":
			len, _ := leb128.ReadVarUint32(r)
			str := make([]byte, len)
			r.Read(str)

			result[name] = string(str)
			break
		}
	}
	return result
}

// LOGIN_SERVER:
// "status:B, inner_id:I, balance:I",
func NewLoginServerPacket(r io.Reader) map[string]interface{} {
	f := getServerFormat(LOGIN_SERVER)
	return CreateServerPacket(LOGIN_SERVER, f, r)
}

// INFO_SERVER:
// "skip:I, skip:I, skip:I, net:L, type:B, name:S, sex:B, tag:B, referrer:I, bday:I, photo:S, photo_byte:B, profile:S",
func NewInfoServerPacket(r io.Reader) map[string]interface{} {
	f := getServerFormat(INFO_SERVER)
	return CreateServerPacket(INFO_SERVER, f, r)
}

// BALANCE_SERVER:
// "bottles:I, reason:B",
func NewBalanceServerPacket(r io.Reader) map[string]interface{} {
	f := getServerFormat(BALANCE_SERVER)
	return CreateServerPacket(BALANCE_SERVER, f, r)
}

// GAME_REWARDS_SERVER:
// "skip:I, id:I, count:I, json:S",
func NewGameRewardsServerPacket(r io.Reader) map[string]interface{} {
	f := getServerFormat(GAME_REWARDS_SERVER)
	return CreateServerPacket(GAME_REWARDS_SERVER, f, r)
}

// BONUS_SERVER:
// "can_collect:B, day:B",
func NewBonusServerPacket(r io.Reader) map[string]interface{} {
	f := getServerFormat(BONUS_SERVER)
	return CreateServerPacket(BONUS_SERVER, f, r)
}

type SERVER_AUTH_RESULT byte

const (
	LOGIN_SUCCESS       SERVER_AUTH_RESULT = 0x00
	LOGIN_FAILED        SERVER_AUTH_RESULT = 0x01
	LOGIN_EXIST         SERVER_AUTH_RESULT = 0x02
	LOGIN_BLOCKED       SERVER_AUTH_RESULT = 0x03
	LOGIN_WRONG_VERSION SERVER_AUTH_RESULT = 0x04
	LOGIN_NO_SEX        SERVER_AUTH_RESULT = 0x05
	LOGIN_CAPTCHA       SERVER_AUTH_RESULT = 0x06

	LOGIN_WAIT_AUTHORIZATION SERVER_AUTH_RESULT = 0xfc
	LOGIN_ERROR              SERVER_AUTH_RESULT = 0xff
)

func (s SERVER_AUTH_RESULT) ToString() string {
	switch s {
	case LOGIN_SUCCESS:
		return "SUCCESS"
	case LOGIN_EXIST:
		return "EXIST"
	case LOGIN_BLOCKED:
		return "BLOCKED"
	case LOGIN_FAILED:
		return "FAILED"
	case LOGIN_CAPTCHA:
		return "CAPTCHA"
	default:
		return "XZ"
	}
}
