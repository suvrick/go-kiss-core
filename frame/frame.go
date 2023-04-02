// Package parser ...
// Пакет для парсинга фреймов в структуру LoginParams
package frame

import (
	"fmt"
	"hash/fnv"
	"net/url"
	"strconv"
	"strings"

	"github.com/suvrick/go-kiss-core/types"
)

const (
	VK byte = 0
	MM byte = 1
	OK byte = 4
	FS byte = 30
	SA byte = 32
	GS byte = 41
	NN byte = 255
)

// Названия ключей для каждого типа соц.сети
type Key struct {
	FrameType byte   `json:"frame_type"`
	LoginID   string `json:"id"`
	Token     string `json:"token"`
	Token2    string `json:"token2"`
	Tag       string `json:"tag"`
	OAuth     string `json:"oauth"`
}

var KEIS = []Key{
	{
		FrameType: 0,
		LoginID:   "viewer_id",
		Token:     "auth_key",
		Token2:    "access_token",
		Tag:       "",
		OAuth:     "OAuth",
	},
	{
		FrameType: 1,
		LoginID:   "vid",
		Token:     "authentication_key",
		Token2:    "session_key",
		Tag:       "",
		OAuth:     "OAuth",
	},
	{
		FrameType: 1,
		LoginID:   "vid",
		Token:     "authentication_key",
		Token2:    "access_token",
		Tag:       "",
		OAuth:     "OAuth",
	},
	{
		FrameType: 4,
		LoginID:   "logged_user_id",
		Token:     "auth_sig",
		Token2:    "session_key",
		Tag:       "",
		OAuth:     "OAuth",
	},
	{
		FrameType: 30,
		LoginID:   "userId",
		Token:     "authKey",
		Token2:    "",
		Tag:       "fotostrana",
		OAuth:     "OAuth",
	},
	{
		FrameType: 32,
		LoginID:   "userId",
		Token:     "authKey",
		Token2:    "sessionKey",
		Tag:       "",
		OAuth:     "OAuth",
	},
}

func Parse3(input string) ([]interface{}, error) {

	result := make([]any, 0)
	login := Parse2(input)

	if login["error"] != nil {
		return nil, fmt.Errorf("[Parse3] %v", login["error"])
	}

	// data := []interface{}{
	// 	types.L(105345504),
	// 	types.B(30),
	// 	types.B(5),
	// 	types.S("7b0a077a088b9e5169bcfc0bf2ee9ae8"),
	// 	types.B(1),
	// 	types.S("5d5b1908c2bae78eeb199db47fc327ac935ccfbd914a38"),
	// 	types.I(0),
	// 	types.I(0),
	// 	types.B(0),
	// 	types.S(""),
	// 	types.B(0),
	// 	types.S(""),
	// 	types.B(0),
	// 	types.S(""),
	// }

	result = append(result, types.L(login["login_id"].(uint64)))
	result = append(result, types.B(login["frame_type"].(byte)))
	result = append(result, types.B(login["device"].(byte)))
	result = append(result, types.S(login["key"].(string)))
	oauth := types.B(0)
	if login["oauth"].(bool) {
		oauth = types.B(1)
	}
	result = append(result, oauth)
	result = append(result, types.S(login["access_token"].(string)))
	result = append(result, types.I(0))
	result = append(result, types.I(0))
	result = append(result, types.B(0))
	result = append(result, types.S(""))
	result = append(result, types.B(0))
	result = append(result, types.S(""))
	result = append(result, types.B(0))
	result = append(result, types.S(""))

	return result, nil
}

func Parse2(input string) map[string]interface{} {
	return Parse(input, KEIS)
}

func Parse(input string, keis []Key) map[string]interface{} {

	result := make(map[string]interface{})
	result["hash"] = getFrameHash(input)
	result["frame"] = input
	result["frame_type"] = NN
	result["frame_type_name"] = getFrameTypeName(NN)

	q, err := url.ParseQuery(input)
	if err != nil {
		result["error"] = fmt.Sprintf("frame parse fail: %s", err.Error())
		return result
	}

	i, c := -1, 2

	for index, key := range keis {

		var counter int

		if q.Has(key.LoginID) {
			counter++
		}

		if q.Has(key.Token) {
			counter++
		}

		if counter == 2 && len(key.Tag) != 0 && strings.Contains(input, key.Tag) {
			counter++
		}

		if counter >= c {
			i = index
			c = counter
		}
	}

	if i > -1 {
		id, err := strconv.ParseUint(q.Get(keis[i].LoginID), 10, 64)
		if err != nil {
			result["error"] = fmt.Sprintf("frame parse fail: %s", err.Error())
			return result
		}

		result["login_id"] = id
		result["device"] = byte(5)
		result["frame_type"] = byte(keis[i].FrameType)
		result["frame_type_name"] = getFrameTypeName(keis[i].FrameType)
		result["key"] = q.Get(keis[i].Token)
		result["oauth"] = q.Has(keis[i].OAuth)
		result["access_token"] = q.Get(keis[i].Token2)
	} else {
		result["error"] = "frame parse fail: unknown frame type"
	}

	return result
}

// getFrameTypeName возращает строковое представления FrameType.
// Если не удалось определить тип frame, то "nn"
func getFrameTypeName(t byte) string {
	switch t {
	case VK:
		return "vk"
	case MM:
		return "mm"
	case OK:
		return "ok"
	case FS:
		return "fs"
	case SA:
		return "sa"
	case GS:
		return "gs"
	default:
		return "nn"
	}
}

func getFrameHash(s string) uint64 {
	hex := fnv.New64a()
	_, err := hex.Write([]byte(s))
	if err != nil {
		return 0
	}
	return hex.Sum64()
}
