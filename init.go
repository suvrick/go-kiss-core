package gokisscore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-interpreter/wagon/wasm/leb128"
)

type Filler struct {
	URL_FULL    string
	URL_VERTION string
	URL_SCRIPT  string

	body_lines   []string
	game_version string

	ClientPackets map[ClientPacketType]ClientPacket
	ServerPackets map[ServerPacketType]ServerPacket

	c_formats []string
	s_formats []string

	Error error
}

type ClientPacketType int
type ServerPacketType int

type ClientPacket struct {
	Type   ClientPacketType
	Format string
	Name   string
}

type ServerPacket struct {
	Type   ServerPacketType
	Format string
	Name   string
}

func NewFiller() *Filler {
	return &Filler{

		URL_FULL:    "",
		URL_VERTION: "https://inspin.me/version.json",
		URL_SCRIPT:  "https://inspin.me/build/v%s/scripts/main.js",

		game_version: "",

		body_lines: make([]string, 0),

		s_formats: make([]string, 0),
		c_formats: make([]string, 0),

		ServerPackets: make(map[ServerPacketType]ServerPacket),
		ClientPackets: make(map[ClientPacketType]ClientPacket),
	}
}

func (filler *Filler) ParseScript() {

	filler.GetVersion()
	filler.GetBody()

	filler.s_formats = filler.SetFormat("PacketServer.FORMATS=[", "];")
	filler.c_formats = filler.SetFormat("PacketClient.FORMATS=[", "];")

	filler.SetServerType()
	filler.SetClientType()

}

func (filler *Filler) SetFormat(pattern_start, pattern_end string) []string {

	if filler.Error != nil {
		return nil
	}

	defer func() {
		r := recover()
		if r != nil {
			msg := fmt.Sprintf("recover from GetFormat.\n %v\n", r)
			filler.Error = errors.New(msg)
			log.Print(msg)
		}
	}()

	isOpen := false
	formats := make([]string, 0)

	//Начинаем переберать ответ от сервера
	for _, line := range filler.body_lines {

		// Ищим совпадения начало шаблона,
		//выставляем флаг открытия для начало парсинга формата
		if strings.Contains(line, pattern_start) {
			isOpen = true
			continue
		}

		if isOpen {

			// проверяем строку на конец шаблона
			if strings.Contains(line, pattern_end) {
				// выходим из цикла
				break
			}

			format := ""
			is_add_char := false // добавить руну к локальному формату format

			// ищим в строке шаблон формата "BBBSSS"
			for _, char := range line {

				if char == '"' {
					is_add_char = !is_add_char

					if !is_add_char {
						formats = append(formats, format)
						format = ""
					}

					continue
				}

				if is_add_char {
					format += string(char)
				}
			}

		}

	}

	if len(formats) == 0 {
		er := fmt.Sprintf("empty format array by pattern_start: \"%s\", pattern_end: \"%s\"", pattern_start, pattern_end)
		filler.Error = errors.New(er)
		return formats
	}

	log.Printf("set format array len %d by pattern_start: \"%s\", pattern_end: \"%s\"", len(formats), pattern_start, pattern_end)
	return formats
}

func (filler *Filler) SetServerType() {

	if filler.Error != nil {
		return
	}

	defer func() {
		r := recover()
		if r != nil {
			msg := fmt.Sprintf("recover from GetServerType.\n %v\n", r)
			filler.Error = errors.New(msg)
			log.Print(msg)
		}
	}()

	pattern := "ServerPacketType[ServerPacketType["

	for _, line := range filler.body_lines {

		if !strings.Contains(line, pattern) {
			continue
		}

		p1 := "\"]="
		p2 := "]=\""

		i1 := strings.Index(line, p1)
		i2 := strings.Index(line, p2)

		if i1 == -1 || i2 == -1 {
			continue
		}

		i1 += len(p1)

		p_type := line[i1:i2]

		p3 := "]=\""
		p4 := "\";"

		i3 := strings.Index(line, p3)
		i4 := strings.Index(line, p4)

		if i3 == -1 || i4 == -1 {
			continue
		}

		i3 += len(p3)

		p_type_name := line[i3:i4]

		id, err := strconv.Atoi(p_type)
		if err != nil {
			continue
		}

		if id >= len(filler.s_formats) {
			//log.Printf("continue id format (id > len(formats)) %d\n", id)
			continue
		}

		pack := ServerPacket{
			Type:   ServerPacketType(id),
			Name:   p_type_name,
			Format: filler.s_formats[id],
		}

		filler.ServerPackets[pack.Type] = pack
	}

	length := len(filler.ServerPackets)
	if length == 0 {
		filler.Error = errors.New("empty map of server packets")
		return
	}

	log.Printf("set server packets len(%d)\n", length)
}

func (filler *Filler) SetClientType() {

	if filler.Error != nil {
		return
	}

	defer func() {
		r := recover()
		if r != nil {
			msg := fmt.Sprintf("recover from GetClientType.\n %v\n", r)
			filler.Error = errors.New(msg)
			log.Print(msg)
		}
	}()

	pattern := "ClientPacketType[ClientPacketType["

	for _, line := range filler.body_lines {

		if !strings.Contains(line, pattern) {
			continue
		}

		p1 := "\"]="
		p2 := "]=\""

		i1 := strings.Index(line, p1)
		i2 := strings.Index(line, p2)

		if i1 == -1 || i2 == -1 {
			continue
		}

		i1 += len(p1)

		p_type := line[i1:i2]

		p3 := "]=\""
		p4 := "\";"

		i3 := strings.Index(line, p3)
		i4 := strings.Index(line, p4)

		if i3 == -1 || i4 == -1 {
			continue
		}

		i3 += +len(p3)

		p_type_name := line[i3:i4]

		id, err := strconv.Atoi(p_type)
		if err != nil {
			continue
		}

		if id >= len(filler.c_formats) {
			//log.Printf("continue id format (id > len(formats)) %d\n", id)
			continue
		}

		pack := ClientPacket{
			Type:   ClientPacketType(id),
			Name:   p_type_name,
			Format: filler.c_formats[id],
		}

		filler.ClientPackets[pack.Type] = pack
	}

	length := len(filler.ClientPackets)
	if length == 0 {
		filler.Error = errors.New("empty map of client packets")
		return
	}

	log.Printf("set client packets len(%d)\n", length)
}

/*

	GetBody
	get script game by url https://inspin.me/build/v%s/scripts/main.js

*/
func (filler *Filler) GetBody() {

	if filler.Error != nil {
		return
	}

	defer func() {
		r := recover()
		if r != nil {
			msg := fmt.Sprintf("recover from GetBody.\n %v\n", r)
			filler.Error = errors.New(msg)
			log.Print(msg)
		}
	}()

	filler.URL_FULL = fmt.Sprintf(filler.URL_SCRIPT, filler.game_version)
	log.Println(filler.URL_FULL)
	resp, err := http.Get(filler.URL_FULL)

	if err != nil {
		filler.Error = err
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		filler.Error = err
		return
	}

	b := string(body)

	log.Println(b)

	b = strings.ReplaceAll(b, " ", "")

	filler.body_lines = strings.Split(b, "\n")

	if len(filler.body_lines) < 2 {
		filler.Error = errors.New("empty response from game server")
		return
	}

	log.Printf("set body len(%d)\n", len(filler.body_lines))
}

/*

	GetVersion
	get version game by url https://inspin.me/version.json

*/
func (filler *Filler) GetVersion() {

	defer func() {
		r := recover()
		if r != nil {
			msg := fmt.Sprintf("recover from GetVersion.\n %v\n", r)
			filler.Error = errors.New(msg)
			log.Print(msg)
		}
	}()

	resp, err := http.Get(filler.URL_VERTION)

	if err != nil {
		filler.Error = err
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		filler.Error = err
		log.Println(err.Error())
		return
	}

	log.Println(string(body))

	type Version struct {
		Version string `json:"version"`
	}

	v := Version{}

	err = json.Unmarshal(body, &v)

	if err != nil {
		filler.Error = err
		return
	}

	// Проверяем заполнения поля Version
	// если длина 0 то ошибка
	if len(v.Version) == 0 {
		filler.Error = errors.New("error of get version")
		return
	}

	// Проверяем валидность версии
	// должно быть числом
	// иначе ошибка
	_, err = strconv.Atoi(v.Version)
	if err != nil {
		msg := fmt.Sprintf("error get game version. invalid string %s", v.Version)
		filler.Error = errors.New(msg)
		return
	}

	filler.game_version = v.Version
	log.Printf("set game version %s\n", filler.game_version)
}

func (filler *Filler) CreateClientPacket(t ClientPacketType, params ...interface{}) []byte {

	p, ok := filler.ClientPackets[t]
	if !ok {
		log.Printf("error create client packet %s.not found client packet %s\n", p.Name, p.Name)
		return nil
	}

	result := make([]byte, 0)
	result = leb128.AppendUleb128(result, uint64(t)) // packet type
	result = leb128.AppendSleb128(result, int64(4))  // device type

	var index int
	var optional bool
	for _, symbol := range p.Format {

		if symbol == ',' {
			optional = true
			continue
		}

		if symbol == ']' || symbol == '[' {
			if optional {
				return result
			}

			log.Printf("error create client packet %s. type 'array' not impliment yet.\n", p.Name)
			return nil
		}

		if index >= len(params) {

			if optional {
				return result
			}

			log.Printf("error create client packet %s. little params\n", p.Name)
			return nil
		}

		value := params[index]

		switch symbol {
		case 'B', 'I', 'L':

			switch v := value.(type) {
			case int8:
				result = leb128.AppendSleb128(result, int64(value.(int8)))
			case uint8:
				result = leb128.AppendUleb128(result, uint64(value.(uint8)))
			case int:
				result = leb128.AppendSleb128(result, int64(value.(int)))
			case uint:
				result = leb128.AppendUleb128(result, uint64(value.(uint)))
			case int32:
				result = leb128.AppendSleb128(result, int64(value.(int32)))
			case uint32:
				result = leb128.AppendUleb128(result, uint64(value.(uint32)))
			case int64:
				result = leb128.AppendSleb128(result, value.(int64))
			case uint64:
				result = leb128.AppendUleb128(result, value.(uint64))

			case float32:
				result = leb128.AppendSleb128(result, int64(value.(float32)))
			case float64:
				result = leb128.AppendUleb128(result, uint64(value.(float64)))

			default:
				log.Printf("error create client packet %s. bad signature for client packet %s.current type: %T, want type %s\n", p.Name, p.Name, v, string(symbol))
				return nil
			}
		case 'S':

			v, ok := value.(string)

			if !ok {
				log.Printf("error create client packet %s. bad signature for client packet %s.current type: %T, want type %s\n", p.Name, p.Name, reflect.TypeOf(value).Kind(), string(symbol))
				return nil
			}

			if len(v) == 0 {
				result = append(result, 0)
				continue
			}

			result = leb128.AppendUleb128(result, uint64(len(v)))
			result = append(result, []byte(v)...)
			result = append(result, 0)
		}

		index++
	}

	return result
}

func (filler *Filler) CreateServerPacket(t ServerPacketType, r io.Reader) (string, []interface{}) {

	p, ok := filler.ServerPackets[t]
	if !ok {
		log.Printf("error create server packet %s.not found server packet %s\n", p.Name, p.Name)
		return "NOT FOUND", nil
	}

	var index int
	var optional bool
	var is_array_type bool
	var array_len uint32

	result := make([]interface{}, 0)

	for _, symbol := range p.Format {

		if symbol == ',' {
			optional = true
			continue
		}

		if symbol == ']' {

			is_array_type = false
			continue
		}

		if symbol == '[' {
			is_array_type = true
			len, err := leb128.ReadVarUint32(r)
			if err != nil {
				log.Printf("error create server packet %s. type 'array' not impliment yet.\n", p.Name)
				return p.Name, nil
			}

			array_len = len
			log.Printf("server packet %s, array len %d\n", p.Name, array_len)
			continue
		}

		switch symbol {
		case 'B', 'I', 'L':
			value, err := leb128.ReadVarint64(r)
			if err != nil {
				log.Printf("error create server packet %s. invalid type. want type %s\n", p.Name, string(symbol))
				return p.Name, nil
			}

			result = append(result, value)
		case 'S':
			len, err := leb128.ReadVarUint32(r)
			if err != nil {
				log.Printf("error create server packet %s. invalid type. want type %s\n", p.Name, string(symbol))
				return p.Name, nil
			}

			str := make([]byte, len)

			_, err = r.Read(str)
			if err != nil {
				log.Printf("error create server packet %s. invalid type. want type %s\n", p.Name, string(symbol))
				return p.Name, nil
			}

			result = append(result, string(str))
		}

		index++
		_ = optional
		_ = is_array_type
	}

	return p.Name, result
}
