package net

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrEmptyParams = errors.New("empty params")
	ErrEmptyResult = errors.New("empty result")
	ErrBadRequest  = errors.New("bad request")
)

type RegConfig struct {
	Start   []byte
	End     []byte
	Pattern []byte
}

type ParserConfig struct {
	HostPath    string
	VersionPath string
	ScriptPath  string
	Version     string

	ClientFormats RegConfig
	ServerFormats RegConfig
	ServerTypes   RegConfig
	ClientTypes   RegConfig
}

func GetDefaultParserConfig() *ParserConfig {
	return &ParserConfig{
		HostPath:    "https://inspin.me/",
		VersionPath: "version.json",
		ScriptPath:  "workers/connection_worker.js",
		ClientFormats: RegConfig{
			Start:   []byte("PacketServer=o,o.FORMATS=["),
			End:     []byte("}"),
			Pattern: []byte(`"[A-Z,\[\]]*"`),
		},
		ServerFormats: RegConfig{
			Start:   []byte("PacketServer=o,o.FORMATS=["),
			End:     []byte("}"),
			Pattern: []byte(`"[A-Z,\[\]]*"`),
		},
		ClientTypes: RegConfig{
			Start:   []byte(".ClientPacketType=void"),
			End:     []byte("}"),
			Pattern: []byte(`.([A-Z_]*)=([0-9]*)`),
		},
		ServerTypes: RegConfig{
			Start:   []byte(".ServerPacketType=void"),
			End:     []byte("}"),
			Pattern: []byte(`.([A-Z_]*)=([0-9]*)`),
		},
	}
}

type Parser struct {
	ParserConfig
	Error error
}

func NewParser() *Parser {
	return &Parser{
		ParserConfig: *GetDefaultParserConfig(),
	}
}

func (p *Parser) Initialize() {

	p.getVersion()
	if p.Error != nil {
		fmt.Printf("%v\n", p.Error)
		return
	}

	fmt.Printf("set game version: %s\n", p.Version)

	body := p.getBody()
	if p.Error != nil {
		fmt.Printf("%v\n", p.Error)
		return
	}

	fmt.Printf("body length: %d\n", len(body))

	/* FORMATS */

	server_formats := p.getFormats(body, p.ServerFormats.Start, p.ServerFormats.End, p.ServerFormats.Pattern)
	fmt.Printf("get server formats length: %d\n", len(server_formats))

	client_formats := p.getFormats(body, p.ClientFormats.Start, p.ClientFormats.End, p.ClientFormats.Pattern)
	fmt.Printf("get client formats length: %d\n", len(client_formats))

	/* TYPES */

	server_types := p.getTypes(body, p.ServerTypes.Start, p.ServerTypes.End, p.ServerTypes.Pattern)
	fmt.Printf("get server types length: %d\n", len(server_types))

	client_types := p.getTypes(body, p.ClientTypes.Start, p.ClientTypes.End, p.ClientTypes.Pattern)
	fmt.Printf("get client types length: %d\n", len(client_types))

	/* GENERATE */

	servers := p.generateData(server_formats, server_types)
	fmt.Printf("generate servers: %d\n", len(servers))

	clients := p.generateData(client_formats, client_types)
	fmt.Printf("generate clients: %d\n", len(clients))

}

func (parser *Parser) generateData(formats []string, types map[int]string) map[int][2]string {
	result := make(map[int][2]string)
	for id, format := range formats {
		name, ok := types[id]
		if ok {
			result[id] = [2]string{
				name,
				format,
			}
		}
	}

	return result
}

func (parser *Parser) getFormats(body []byte, start []byte, end []byte, pattern []byte) []string {

	body = parser.split(body, start, end)
	if len(body) == 0 {
		return nil
	}

	// fmt.Printf("body: %s", body)

	reg, err := regexp.Compile(string(pattern))
	if err != nil {
		return nil
	}

	finded := reg.FindAll(body, -1)

	result := make([]string, 0)

	for _, v := range finded {
		result = append(result, string(v))
	}

	return result
}

func (parser *Parser) getTypes(body []byte, start []byte, end []byte, pattern []byte) map[int]string {
	body = parser.split(body, start, end)
	if len(body) == 0 {
		return nil
	}

	//fmt.Printf("body types: %s", body)

	reg, err := regexp.Compile(string(pattern))
	if err != nil {
		return nil
	}

	finded := reg.FindAllSubmatch(body, -1)

	result := make(map[int]string, len(finded))

	for _, v := range finded {

		if len(v) != 3 {
			continue
		}

		key := v[2]
		value := v[1]
		id, err := strconv.Atoi(string(key))
		if err == nil {
			result[id] = string(value)
		}
	}

	return result
}

func (parser *Parser) getBody() []byte {
	url := fmt.Sprintf("%s%s%s", parser.HostPath, parser.Version, parser.ScriptPath)
	body, err := parser.request(url)
	if err != nil {
		parser.Error = err
		return nil
	}
	return body
}

func (parser *Parser) getVersion() {

	url := fmt.Sprintf("%s%s", parser.HostPath, parser.VersionPath)

	body, err := parser.request(url)
	if err != nil {
		parser.Error = err
		return
	}

	type version struct {
		V string `json:"browser"`
	}

	v := version{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		parser.Error = err
		return
	}

	if len(v.V) == 0 {
		parser.Error = ErrEmptyResult
		return
	}

	parser.Version = v.V
}

func (parser *Parser) request(url string) ([]byte, error) {

	resp, err := http.Get(url)
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

	// delete spacial symbol
	body = bytes.ReplaceAll(body, []byte{9}, []byte{})  // remove "	"
	body = bytes.ReplaceAll(body, []byte{32}, []byte{}) // remove " "
	body = bytes.ReplaceAll(body, []byte{10}, []byte{}) // remove "\n"
	body = bytes.ReplaceAll(body, []byte{13}, []byte{}) // remove "\r"

	return body, nil
}

func (parser *Parser) split(body []byte, start []byte, end []byte) []byte {
	index_start := bytes.Index(body, start)
	if index_start == -1 {
		return nil
	}

	index_start = index_start + len(start)

	index_end := bytes.Index(body[index_start:], end)
	if index_end == -1 {
		return nil
	}

	index_end = index_start + index_end

	fmt.Printf("index_start: %d, index_end %d\n", index_start, index_end)

	return body[index_start:index_end]
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
