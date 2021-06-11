package packets

import (
	"strings"

	"github.com/go-interpreter/wagon/wasm/leb128"
)

type ClientPacketType uint16

const (
	// net_id:I, type:B, device:B, auth_key:S, oauth:B, session_key:S, referrer:I, tag:I, appicationID:B, timestamp:S, language:B, utm_source:S, sex:B, captcha:S
	LOGIN ClientPacketType = 4
	// good_id:I, cost:I, target_id:I, data:I, price_type:B, count:I, hash:S, params: S
	BUY ClientPacketType = 6
	// reward_id:I
	GAME_REWARDS_GET ClientPacketType = 11
	// inner_id:I
	VIEW ClientPacketType = 17
	// ""
	BONUS ClientPacketType = 61
	// ""
	GET_BDAY_REWARD ClientPacketType = 163
	// type:B
	PROFILE_REWARD_GET ClientPacketType = 170
	// ""
	ROULETTE_ROLL ClientPacketType = 252
	// player_id:I
	GET_ADMIRE_BONUS ClientPacketType = 266
	// type:I, data:I, count:I
	COUNT ClientPacketType = 20
)

var c_format = map[ClientPacketType]string{
	LOGIN:              "net_id:I, type:B, device:B, auth_key:S, oauth:B, session_key:S, referrer:I, tag:I, appicationID:B, timestamp:S, language:B, utm_source:S, sex:B",
	BUY:                "good_id:I, cost:I, target_id:I, data:I, price_type:B, count:I, hash:S, params: S",
	GAME_REWARDS_GET:   "reward_id:I",
	VIEW:               " inner_id:I",
	BONUS:              "",
	GET_BDAY_REWARD:    "",
	PROFILE_REWARD_GET: "type:B",
	ROULETTE_ROLL:      "",
	GET_ADMIRE_BONUS:   "player_id:I",
	COUNT:              "type:I, data:I, count:I",
}

func getClientFormat(t ClientPacketType) string {

	line, ok := c_format[t]

	if !ok {
		return ""
	}

	return line
}

func CreateClientPacket(t ClientPacketType, format string, data map[string]interface{}) []byte {

	result := make([]byte, 0)

	result = leb128.AppendUleb128(result, uint64(t)) // packet type
	result = leb128.AppendSleb128(result, int64(4))  // device type

	format = strings.ReplaceAll(format, " ", "")
	words := strings.Split(format, ",")

	if words == nil {
		return result
	}

	for _, word := range words {

		s := strings.Split(word, ":")

		if len(s) != 2 {
			continue
		}

		name := s[0]
		code := s[1]

		value := data[name]

		if value == nil {
			result = append(result, 0)
			continue
		}

		switch code {
		case "B", "I", "L":

			switch value.(type) {
			case uint8:
				result = leb128.AppendUleb128(result, uint64(value.(uint8)))
				break
			case int16:
				result = leb128.AppendUleb128(result, uint64(value.(int16)))
				break
			case int:
				result = leb128.AppendUleb128(result, uint64(value.(int)))
				break
			case int32:
				result = leb128.AppendUleb128(result, uint64(value.(int32)))
				break
			case uint32:
				result = leb128.AppendUleb128(result, uint64(value.(uint32)))
				break
			case uint64:
				result = leb128.AppendUleb128(result, value.(uint64))
				break
			}
			break
		case "S":

			v, ok := value.(string)

			if !ok {
				result = append(result, 0)
				continue
			}

			if len(v) == 0 {
				result = append(result, 0)
				continue
			}

			result = leb128.AppendUleb128(result, uint64(len(v)))
			result = append(result, []byte(v)...)
			break
		}
	}

	return result
}

// LOGIN
// "net_id:L, type:B, device:B, auth_key:S, oauth:B, session_key:S
func NewLoginPacketClient(lp map[string]interface{}) []byte {
	f := getClientFormat(LOGIN)
	return CreateClientPacket(LOGIN, f, lp)
}

// BUY
// "good_id:I, cost:I, target_id:I, data:I, price_type:B, count:I, hash:S, params: S"
func NewBuyPacketClient(good_id, cost, target_id, data int, price_type byte, count int, hash, params string) []byte {
	f := getClientFormat(BUY)
	return CreateClientPacket(BUY, f, map[string]interface{}{
		"good_id":    good_id,
		"cost":       cost,
		"target_id":  target_id,
		"data":       data,
		"price_type": price_type,
		"count":      count,
		"hash":       hash,
		"params":     params,
	})
}

// GAME_REWARDS_GET
// "reward_id:I"
func NewGameRewardsGetPacketClient(reward_id interface{}) []byte {
	f := getClientFormat(GAME_REWARDS_GET)
	return CreateClientPacket(GAME_REWARDS_GET, f, map[string]interface{}{
		"reward_id": reward_id,
	})
}

// VIEW
// "inner_id:I",
func NewViewPacketClient(inner_id int) []byte {
	f := getClientFormat(VIEW)
	return CreateClientPacket(VIEW, f, map[string]interface{}{
		"inner_id": inner_id,
	})
}

// BONUS
// ""
func NewBonusPacketClient() []byte {
	f := getClientFormat(BONUS)
	return CreateClientPacket(BONUS, f, nil)
}

// GET_BDAY_REWARD
// ""
func NewGetBDayRewardPacketClient() []byte {
	f := getClientFormat(GET_BDAY_REWARD)
	return CreateClientPacket(GET_BDAY_REWARD, f, nil)
}
