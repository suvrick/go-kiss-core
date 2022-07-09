package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/suvrick/go-kiss-core/packets"
)

var (
	ErrEmptyParams = errors.New("empty params")
	ErrEmptyResult = errors.New("empty result")
	ErrBadRequest  = errors.New("bad request")
)

type Parser struct {
	game_host           string
	name_version_script string
	name_game_script    string
	name_worker_script  string
	query_version       string

	Error error
}

/*

	Доделать иницилизацию с конфигом
	NewParser(config *Config)

*/
func NewParser() *Parser {
	return &Parser{
		game_host:           "https://inspin.me/",
		name_version_script: "version.json",
		name_game_script:    "scripts/main.js",
		name_worker_script:  "workers/connection_worker.js",
		query_version:       "",
	}
}

//https://inspin.me/build/v2100/workers/connection_worker.js

func (p *Parser) Initialize() {

	init_rewards()

	version, err := getVersion(p.game_host, p.name_version_script)

	if err != nil {
		p.Error = errors.New(fmt.Sprint("error [getVersion] > ", err))
		return
	}

	if version == p.query_version {
		return
	}

	p.query_version = version

	body_lines, err := getBody(p.game_host, p.query_version, p.name_game_script)
	if err != nil {
		p.Error = errors.New(fmt.Sprint("error [getBody] > ", err))
		return
	}

	c_packets := p.InitClientDict(body_lines)
	packets.SetClientPakets(&c_packets)
	if p.Error != nil {
		p.Error = errors.New(fmt.Sprint("error [InitClientDict] > ", err))
		return
	}

	s_packets := p.InitServerDict(body_lines)
	packets.SetServerPacket(&s_packets)
	if p.Error != nil {
		p.Error = errors.New(fmt.Sprint("error [InitServerDict] > ", err))
		return
	}
}

func (p *Parser) InitClientDict(body_lines []string) map[uint64]packets.Packet {

	result := make(map[uint64]packets.Packet)

	c_format, err := setFormat(body_lines, "PacketClient.FORMATS=[", "];")
	if err != nil {
		p.Error = errors.New(fmt.Sprint("error [setFormat] > ", err))
		return nil
	}

	c_types, err := setType(body_lines, "ClientPacketType[ClientPacketType[")
	if err != nil {
		p.Error = errors.New(fmt.Sprint("error [setType] > ", err))
		return nil
	}

	//Заполнения клиенских пакетов
	for id, name := range c_types {

		if id >= uint64(len(c_format)) {
			continue
		}

		format := c_format[id]
		result[id] = packets.Packet{
			Name:   name,
			Type:   id,
			Format: format,
		}
	}

	return result
}

func (p *Parser) InitServerDict(body_lines []string) map[uint64]packets.Packet {
	result := make(map[uint64]packets.Packet)
	s_format, err := setFormat(body_lines, "PacketServer.FORMATS=[", "];")
	if err != nil {
		p.Error = errors.New(fmt.Sprint("error [setFormat] > ", err))
		return nil
	}

	s_types, err := setType(body_lines, "ServerPacketType[ServerPacketType[")
	if err != nil {
		p.Error = errors.New(fmt.Sprint("error [setType] > ", err))
		return nil
	}

	for id, name := range s_types {

		if id >= uint64(len(s_format)) {
			continue
		}

		format := s_format[id]
		result[id] = packets.Packet{
			Name:   name,
			Type:   id,
			Format: format,
		}
	}

	return result
}

/*
	setFormat

	body_line - строковый массив скрипта игры

	pattern_start - строка вхождения типа "PacketServer.FORMATS=[". Старт запуска внутренего цикла персинга результирующего массива

	pattern_end - точка выхода, конечная граница парсинга "];"

	return возращаем массив форматов [ "S", "B,IIBI[B]IIIISBBIBS" ...] и ошибку

*/
func setFormat(body_lines []string, pattern_start, pattern_end string) ([]string, error) {

	isOpen := false
	formats := make([]string, 0)

	//Начинаем переберать ответ от сервера
	for _, line := range body_lines {

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
		err := fmt.Sprintf("empty result array. pattern_start: \"%s\", pattern_end: \"%s\"", pattern_start, pattern_end)
		return formats, errors.New(err)
	}

	return formats, nil
}

/*
	setType

	body_line - строковый массив скрипта игры

	pattern - строка вхождения  "ClientPacketType[ClientPacketType["

	return возращаем словарь вида [ {4:"LOGIN"}, {7:"BALANCE"}... ] и ошибку

*/
func setType(body_lines []string, pattern string) (map[uint64]string, error) {

	result := make(map[uint64]string)

	for _, line := range body_lines {

		if !strings.Contains(line, pattern) {
			continue
		}

		var name strings.Builder
		var id strings.Builder

		is_add_rune := false
		is_add_num := false

		for _, r := range line {

			if r == '"' {
				is_add_rune = !is_add_rune
				continue
			}

			if is_add_rune {
				name.WriteRune(r)
				continue
			}

			if r == '=' {
				is_add_num = true
				continue
			}

			if is_add_num {
				if unicode.IsDigit(r) {
					id.WriteRune(r)
					continue
				}

				num, _ := strconv.Atoi(id.String())
				result[uint64(num)] = name.String()
				break
			}
		}
	}

	if len(result) == 0 {
		err := fmt.Sprintf("empty result array. pattern: \"%s\"", pattern)
		return result, errors.New(err)
	}

	return result, nil
}

/*

	getBody

	host_url - "https://inspin.me/"

	query_version - "build/v1790/"

	name_script - "scripts/main.js"

	return []string script game and error

*/
func getBody(host_url string, query_version string, name_script string) ([]string, error) {

	if len(host_url) == 0 || len(query_version) == 0 || len(name_script) == 0 {
		return nil, ErrEmptyParams
	}

	full_url := fmt.Sprintf("%s%s%s", host_url, query_version, name_script)

	resp, err := http.Get(full_url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrBadRequest
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResult
	}

	body_string := string(body)
	body_string = strings.ReplaceAll(body_string, " ", "")
	body_lines := strings.Split(body_string, "\n")

	return body_lines, nil
}

func (p *Parser) GetBoby() ([]string, error) {
	lines, err := getBody(p.game_host, p.query_version, p.name_game_script)
	return lines, err
}

/*

	getVersion

	url - "https://inspin.me/"

	name_script - "version.json"

	return string of format "build/v1774/" and error

*/
func getVersion(url string, script_name string) (string, error) {

	if len(url) == 0 || len(script_name) == 0 {
		return "", ErrEmptyParams
	}

	full_url := fmt.Sprintf("%s%s", url, script_name)

	resp, err := http.Get(full_url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", ErrBadRequest
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	type version struct {
		Version string `json:"browser"`
	}

	v := version{}
	err = json.Unmarshal(body, &v)

	if err != nil {
		return "", err
	}

	if len(v.Version) == 0 {
		return "", ErrEmptyResult
	}

	return v.Version, nil
}

func (p *Parser) GetVersion() (string, error) {
	v, err := getVersion(p.game_host, p.name_version_script)
	p.query_version = v
	return v, err
}

var rewards []string

type Caption struct {
	Ru string `json:"ru"`
	En string `json:"en"`
}

type Obj struct {
	ID       uint64  `json:"id"`
	Captions Caption `json:"captions"`
}

func GetReward(rewardID uint64) string {
	pattern := fmt.Sprintf("\"id\": %d,", rewardID)
	for _, r := range rewards {
		if strings.Index(r, pattern) > 0 {
			return r
		}
	}

	return ""
}

func init_rewards() error {

	resp, err := http.Get("https://bottleconf.realcdn.ru/rewards.json")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ErrBadRequest
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	jsons := strings.ReplaceAll(string(body), "\t", "")
	lines := strings.Split(jsons, "\n")

	rewards = make([]string, 0)

	for _, line := range lines {
		if strings.Index(line, "id") > 0 && strings.Index(line, "captions") > 0 && strings.Index(line, "content") > 0 {

			if line[len(line)-1] == ',' {
				line = line[:len(line)-1]
			}

			rewards = append(rewards, line)
		}
	}

	return nil
}
