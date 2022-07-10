package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

var instance *Parser
var once sync.Once

type Parser struct {
	ParserConfig
	Error error

	clients map[int][2]string
	servers map[int][2]string
}

func NewParser(config *ParserConfig) *Parser {

	if config == nil {
		config = GetDefaultParserConfig()
	}

	if instance == nil {
		instance = &Parser{
			ParserConfig: *config,
		}
		once.Do(instance.initialize)
	}

	return instance
}

func (p *Parser) initialize() {

	p.getVersion()
	if p.Error != nil {
		return
	}

	fmt.Printf("set game version: %s\n", p.Version)

	body := p.getBody()
	if p.Error != nil {
		return
	}

	//fmt.Printf("body length: %d\n", len(body))

	p.parseBody(body)
}

func (parser *Parser) GetClientsMeta() map[int][2]string {
	return parser.clients
}

func (parser *Parser) GetServersMeta() map[int][2]string {
	return parser.servers
}

func (parser *Parser) parseBody(body []byte) {
	/* FORMATS */

	server_formats := parser.getFormats(body, parser.ServerFormats.Start, parser.ServerFormats.End, parser.ServerFormats.Pattern)
	if len(server_formats) == 0 {
		parser.Error = fmt.Errorf("[net.getFormats] getFormats by server return zero len length. start: %#v, end: %#v, pattern: %#v", string(parser.ServerFormats.Start), string(parser.ServerFormats.End), string(parser.ServerFormats.Pattern))
		return
	}

	client_formats := parser.getFormats(body, parser.ClientFormats.Start, parser.ClientFormats.End, parser.ClientFormats.Pattern)
	if len(client_formats) == 0 {
		parser.Error = fmt.Errorf("[net.getFormats] getFormats by client return zero length. start %s, end: %s, pattern: %s", parser.ClientFormats.Start, parser.ClientFormats.End, parser.ClientFormats.Pattern)
		return
	}

	/* TYPES */

	server_types := parser.getTypes(body, parser.ServerTypes.Start, parser.ServerTypes.End, parser.ServerTypes.Pattern)
	if len(server_types) == 0 {
		parser.Error = fmt.Errorf("[net.getTypes] getTypes by server return zero length. start %s, end: %s, pattern: %s", parser.ServerTypes.Start, parser.ServerTypes.End, parser.ServerTypes.Pattern)
		return
	}

	client_types := parser.getTypes(body, parser.ClientTypes.Start, parser.ClientTypes.End, parser.ClientTypes.Pattern)
	if len(client_types) == 0 {
		parser.Error = fmt.Errorf("[net.getTypes] getTypes by client return zero length. start %s, end: %s, pattern: %s", parser.ClientTypes.Start, parser.ClientTypes.End, parser.ClientTypes.Pattern)
		return
	}

	/* GENERATE */

	parser.servers = parser.generateData(server_formats, server_types)
	if len(parser.servers) == 0 {
		parser.Error = fmt.Errorf("[net.generateData] generateData by server return zero length. formats len: %d, types len: %d", len(server_formats), len(server_types))
		return
	}

	parser.clients = parser.generateData(client_formats, client_types)
	if len(parser.clients) == 0 {
		parser.Error = fmt.Errorf("[net.generateData] generateData by client return zero length. formats len: %d, types len: %d", len(client_formats), len(client_types))
		return
	}

	fmt.Printf("parse success! servers: %d, clients: %d\n", len(parser.servers), len(parser.clients))
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

	finded := reg.FindAllSubmatch(body, -1)

	result := make([]string, 0)

	for _, v := range finded {
		if len(v) == 2 {
			result = append(result, string(v[1]))
		}
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
		parser.Error = fmt.Errorf("[net.getBody] fail get body. %s", err)
		return nil
	}
	return body
}

func (parser *Parser) getVersion() {

	url := fmt.Sprintf("%s%s", parser.HostPath, parser.VersionPath)

	body, err := parser.request(url)
	if err != nil {
		parser.Error = fmt.Errorf("[net.getVersion] fail get version. %s", err)
		return
	}

	type version struct {
		V string `json:"browser"`
	}

	v := version{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		parser.Error = fmt.Errorf("[net.getVersion] fail get version. %s", err)
		return
	}

	if len(v.V) == 0 {
		parser.Error = fmt.Errorf("[net.getVersion] fail get version. %s", ErrEmptyResult)
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
		return nil, fmt.Errorf("%s. request url: %s", ErrBadRequest, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("%s. request url: %s", ErrEmptyResult, url)
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

	//fmt.Printf("index_start: %d, index_end %d\n", index_start, index_end)

	return body[index_start:index_end]
}
