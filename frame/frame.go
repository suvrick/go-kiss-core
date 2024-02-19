// Package parser ...
// Пакет для парсинга фреймов в структуру LoginParams
package frame

import (
	"bytes"
	"encoding/json"
	"errors"
	"hash/fnv"
	"net/url"
	"strings"
	"sync"
)

func init() {
	data, _ := Unmarshal([]byte(keys_json))
	f := New()
	f.Initialize(data)
}

var instance IFrameManager
var once sync.Once

type IFrameManager interface {
	Initialize(words []t_words) error
	Parse(input string) (FrameDTO, error)
}

type Frame struct {
	keys []t_words
}

func (f *Frame) Initialize(data []t_words) error {
	f.keys = make([]t_words, len(data))
	copy(f.keys, data)
	return nil
}

func (f *Frame) Parse(input string) (FrameDTO, error) {
	result := FrameDTO{}

	if f.keys == nil {
		return result, ErrFrameParserNotInit
	}

	if len(f.keys) == 0 {
		return result, ErrFrameParserEmptyKeys
	}

	input = strings.TrimSpace(input)

	// TODO: поискать способ удаления сразу всех спец.символов
	input = strings.Replace(input, "\r", "", -1)

	if len(input) == 0 {
		return result, ErrEmptyString
	}

	q, err := url.ParseQuery(input)
	if err != nil {
		return result, ErrInvalidFrame
	}

	// TODO: возвращаем первое совпадение всех ключей. Хм-м что-то как-то стрёмно
	for _, p := range f.keys {

		if q.Has(p.LoginID) && q.Has(p.Token) && q.Has(p.Token2) {

			id := q.Get(p.LoginID)

			// TODO: костыль для фотостраны
			if strings.Index(input, "fotostrana") > 0 {
				p.FrameType = 30
			}

			// result.Hash = getHex(input)
			result.ID = id
			result.NetType = p.FrameType
			result.Key = q.Get(p.Token)
			result.OAuth = 0
			result.AccessToken = q.Get(p.Token2)
			result.StringField = ""

			return result, nil
		}
	}

	return result, ErrFrameTypeNotFound
}

func New() IFrameManager {
	if instance == nil {
		once.Do(func() {
			instance = &Frame{
				keys: make([]t_words, 0),
			}
		})
	}
	return instance
}

const (
	vk byte = 0
	mm byte = 1
	ok byte = 4
	fs byte = 30
	gs byte = 41
	sa byte = 32
	nn byte = 255
)

type FrameDTO struct {
	ID          string
	NetType     byte
	Key         string
	OAuth       byte
	AccessToken string
	StringField string
}

// Названия ключей для каждого типа соц.сети
// config.json
type t_words struct {
	FrameType byte   `json:"frame_type"`
	LoginID   string `json:"id"`
	Token     string `json:"token"`
	Token2    string `json:"token2"`
}

func Unmarshal(data []byte) ([]t_words, error) {
	result := make([]t_words, 0)
	decoder := json.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&result)
	return result, err
}

/*
GetFrameTypeName возращает строковое представления f_type.

Если не удалось определить тип frame, то "nn"
*/
func GetFrameTypeName(t byte) string {
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

const keys_json = `
[
    {
        "frame_type": 0,
        "id": "viewer_id",
        "token": "auth_key",
        "token2": "access_token"
    },
    {
        "frame_type": 1,
        "id": "vid",
        "token": "authentication_key",
        "token2": "session_key"
    },
    {
        "frame_type": 1,
        "id": "vid",
        "token": "authentication_key",
        "token2": "access_token"
    },
    {
        "frame_type": 4,
        "id": "logged_user_id",
        "token": "auth_sig",
        "token2": "session_key"
    },
    {
        "frame_type": 30,
        "id": "userId",
        "token": "authKey",
        "token2": "fotostrana"
    },
    {
        "frame_type": 32,
        "id": "userId",
        "token": "authKey",
        "token2": "sessionKey"
    }
]
`
