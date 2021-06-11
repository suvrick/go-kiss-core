// Package parser ...
// Пакет для парсинга фреймов в структуру LoginParams
package parser

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// ParseParamKey ..
type ParseParamKey struct {
	SType     string `json:"stype"`
	STypeCode int16  `json:"stypecode"`
	ID        string `json:"id"`
	Token     string `json:"token"`
	Token2    string `json:"token2"`
}

const (
	debug = false
)

// Шаблоны для распарсивания URL
var parseParams []ParseParamKey

//go:embed config.json
var pathParserConfig []byte

// init ...
func init() {
	Initialize()
}

// Initialize ...
// Загрузка ключей для парсинга строки
func Initialize() error {

	parseParams = []ParseParamKey{}
	if err := json.Unmarshal(pathParserConfig, &parseParams); err != nil {
		log.Fatalf("Parser [init] >> %s", err.Error())
		return err
	}

	return nil
}

// NewLoginParams ...
//
// Парсим url в map[string]interface структуру
//
// socialCode - тип регистрации (int16)
//
// Если неудача -> nil
//
func NewLoginParams(input string) map[string]interface{} {

	// Удаляем пробелы и спец.символы
	input = strings.TrimSpace(input)
	input = strings.Replace(input, "\r", "", -1)

	if len(input) == 0 {
		return nil
	}

	// Пытаемся получить map[ключ]=значения, query элементов URL
	query, err := url.ParseQuery(input)
	if err != nil {
		return nil
	}

	/*
		Тут уже начинается хардкор :D
		Пытаемся как-нибудь определить тип социалки
	*/

	for _, p := range parseParams {
		if strings.Contains(input, p.ID) &&
			strings.Contains(input, p.Token) &&
			strings.Contains(input, p.Token2) {

			strID := query.Get(p.ID)
			loginID, err := strconv.ParseUint(strID, 10, 64)
			if err != nil || loginID == 0 {
				return nil
			}

			data := map[string]interface{}{
				"net_id":      loginID,
				"type":        p.STypeCode,
				"type_name":   p.SType,
				"device":      byte(4),
				"auth_key":    query.Get(p.Token),
				"oauth":       1,
				"session_key": query.Get(p.Token2),
				"uid":         fmt.Sprintf("%s%d", p.SType, loginID),
			}

			return data
		}
	}

	// Возращаем дефолт 0xFF
	return nil
}
