// Package parser ...
// Пакет для парсинга фреймов в структуру LoginParams
package frame

import (
	"hash/fnv"
	"net/url"
	"strconv"
	"strings"
)

const (
	VK int = 0
	MM int = 1
	OK int = 4
	FS int = 30
	SA int = 32
	GS int = 41
	NN int = 255
)

// Названия ключей для каждого типа соц.сети
type t_words struct {
	FrameType int    `json:"frame_type"`
	LoginID   string `json:"id"`
	Token     string `json:"token"`
	Token2    string `json:"token2"`
	Tag       string `json:"tag"`
	OAuth     string `json:"oauth"`
}

var QUERIES = []t_words{
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

func Parse2(input string) map[string]interface{} {
	return Parse(input, QUERIES)
}

func Parse(input string, words []t_words) map[string]interface{} {

	result := make(map[string]interface{})
	result["id"] = getHex(input)
	result["frame"] = input
	result["frame_type"] = NN
	result["frame_type_name"] = getFrameTypeName(NN)

	q, err := url.ParseQuery(input)
	if err != nil {
		return result
	}

	i, c := -1, 2

	for index, word := range words {

		var counter int

		if q.Has(word.LoginID) {
			counter++
		}

		if q.Has(word.Token) {
			counter++
		}

		if counter == 2 && len(word.Tag) != 0 && strings.Contains(input, word.Tag) {
			counter++
		}

		if counter >= c {
			i = index
			c = counter
		}
	}

	if i > -1 {
		result["login_id"], result["error"] = strconv.ParseUint(q.Get(words[i].LoginID), 10, 64)
		result["device"] = 5
		result["frame_type"] = words[i].FrameType
		result["frame_type_name"] = getFrameTypeName(words[i].FrameType)
		result["key"] = q.Get(words[i].Token)
		result["oauth"] = q.Has(words[i].OAuth)
		result["access_token"] = q.Get(words[i].Token2)
	}

	return result
}

// getFrameTypeName возращает строковое представления FrameType.
// Если не удалось определить тип frame, то "nn"
func getFrameTypeName(t int) string {
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

func getHex(s string) uint64 {
	hex := fnv.New64a()
	_, err := hex.Write([]byte(s))
	if err != nil {
		return 0
	}
	return hex.Sum64()
}
