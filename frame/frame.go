// Package parser ...
// Пакет для парсинга фреймов в структуру LoginParams
package frame

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"hash/fnv"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/suvrick/go-kiss-core/packets/client"
)

//go:embed config.json
var frameConfigBytes []byte

var (
	//Пустая строка
	ErrEmptyString = errors.New("frame parse error.empty string")
	//Ошибка при разборе строки в map[string]string
	ErrInvalidFrame = errors.New("frame parse error.invalid frame")
	//Нет соотвествий по шаблону.Не смог определить тип социальной сети.
	ErrFrameTypeNotFound = errors.New("frame parse error.frame type not found")
	//Не иницилизирован словарь шаблонов.Вызовите Initialize(...)
	ErrFrameParserNotInit = errors.New("frame parse error.not initialize")
	//Не иницилизирован словарь шаблонов.Вызовите Initialize(...)
	ErrFrameParserEmptyKeys = errors.New("frame parse error. keys is empty")
	//Не смог конвертировать LoginID из строки в uint64
	ErrQueryParametrMiss = errors.New("frame parse error.invalid loginID")
)

type Frame struct {
	log  *log.Logger
	keys []t_words
	Err  error
}

const (
	vk uint16 = 0
	mm uint16 = 1
	ok uint16 = 4
	fs uint16 = 30
	gs uint16 = 41
	sa uint16 = 32
	nn uint16 = 255
)

// Названия ключей для каждого типа соц.сети
// config.json
type t_words struct {
	FrameType uint16 `json:"frame_type"`
	LoginID   string `json:"id"`
	Token     string `json:"token"`
	Token2    string `json:"token2"`
}

var instance *Frame = nil
var once sync.Once

func NewFrameDefault() *Frame {
	return NewFrame(log.Default())
}

// Singleton
func NewFrame(logger *log.Logger) *Frame {

	once.Do(func() {
		instance = &Frame{
			keys: make([]t_words, 0),
			log:  logger,
		}

		if logger == nil {
			instance.log = log.Default()
		}

		instance.load_keys()
	})

	return instance
}

func (f *Frame) load_keys() {
	decoder := json.NewDecoder(bytes.NewReader(frameConfigBytes))
	err := decoder.Decode(&f.keys)
	if err != nil {
		f.Err = err
		return
	}
}

/*
Функция для парсинга строки в словарь элементов.
Принимает на вход строку frame
Возращает словарь с определенными ключами

result: [

		"bot_id": ""
		"frame_type": 32
		"login_id": 114941701
		"token": "33513e2ce85cabfd6ec59d827aa28cea"
		"token2": "67f5e4f7a90144c5eba1b91694132904"
	],

error: nil
*/
func (f *Frame) Parse(input string) (map[string]interface{}, error) {

	result := make(map[string]interface{}, 5)

	if f.keys == nil {
		return result, ErrFrameParserNotInit
	}

	if len(f.keys) == 0 {
		return result, ErrFrameParserEmptyKeys
	}

	result["frame_type"] = nn

	input = strings.TrimSpace(input)
	input = strings.Replace(input, "\r", "", -1)

	if len(input) == 0 {
		return result, ErrEmptyString
	}

	result["bot_id"] = getHex(input)

	queries, err := url.ParseQuery(input)
	if err != nil {
		return result, ErrInvalidFrame
	}

	for _, p := range f.keys {

		if queries.Has(p.LoginID) && queries.Has(p.Token) && queries.Has(p.Token2) {

			id, err := strconv.ParseUint(queries.Get(p.LoginID), 10, 64)
			if err != nil {
				return result, ErrQueryParametrMiss
			}

			if strings.Index(input, "fotostrana") > 0 {
				p.FrameType = 30
			}

			result["login_id"] = id
			result["frame_type"] = p.FrameType
			result["token"] = queries.Get(p.Token)
			result["token2"] = queries.Get(p.Token2)

			return result, nil
		}
	}

	return result, ErrFrameTypeNotFound
}

func (f *Frame) Parse2(data []byte) (*client.Login, error) {

	input := string(data)

	if f.keys == nil {
		return nil, ErrFrameParserNotInit
	}

	if len(f.keys) == 0 {
		return nil, ErrFrameParserEmptyKeys
	}

	input = strings.TrimSpace(input)
	input = strings.Replace(input, "\r", "", -1)

	if len(input) == 0 {
		return nil, ErrEmptyString
	}

	queries, err := url.ParseQuery(input)
	if err != nil {
		return nil, ErrInvalidFrame
	}

	for _, p := range f.keys {

		if queries.Has(p.LoginID) && queries.Has(p.Token) && queries.Has(p.Token2) {

			id, err := strconv.ParseUint(queries.Get(p.LoginID), 10, 64)
			if err != nil {
				return nil, ErrQueryParametrMiss
			}

			if strings.Index(input, "fotostrana") > 0 {
				p.FrameType = 30
			}

			return &client.Login{
				ID:          id,
				NetType:     p.FrameType,
				DeviceType:  5,
				Key:         queries.Get(p.Token),
				AccessToken: queries.Get(p.Token2),
			}, nil
		}
	}

	return nil, ErrFrameTypeNotFound
}

func (f *Frame) GetValue(intput string) (uint32, []interface{}, error) {
	//{113594657, 32, 4, "7a2b140e7b42935768c040a54ade4cfc", 0, "8c9991f3e49ef7d20d33432d1534e378"}
	r, e := f.Parse(intput)
	return r["bot_id"].(uint32), []interface{}{
		r["login_id"].(uint64),
		r["frame_type"].(uint16),
		uint16(4),
		r["token"].(string),
		uint8(0),
		r["token2"].(string),
	}, e
}

// GetFrameTypeName возращает строковое представления f_type.Если не удалось определить тип frame, то "nn"
func GetFrameTypeName(t uint16) string {
	switch t {
	case vk:
		return "vk"
	case mm:
		return "mm"
	case ok:
		return "ok"
	case fs:
		return "fs"
	case sa:
		return "sa"
	case gs:
		return "gs"
	default:
		return "nn"
	}
}

func getHex(s string) uint32 {

	hex := fnv.New32a()

	_, err := hex.Write([]byte(s))
	if err != nil {
		return 0
	}

	return hex.Sum32()
}
