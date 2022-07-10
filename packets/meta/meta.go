package meta

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

var instance *Meta
var once sync.Once
var configure *MetaConfig

type Meta struct {
	MetaConfig
	Error error

	clients map[int][2]string
	servers map[int][2]string
}

func NewMeta() *Meta {

	if instance == nil {
		if configure == nil {
			configure = GetDefaultMetaConfig()
		}

		instance = &Meta{
			MetaConfig: *configure,
		}
		once.Do(instance.initialize)
	}

	return instance
}

func SetConfig(config *MetaConfig) {
	configure = config
	instance = nil
}

func (p *Meta) initialize() {

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

// return id, name, format, ok
func GetClientMeta(id int) (int, string, string, bool) {

	if instance == nil {
		NewMeta()
	}

	r, ok := instance.clients[id]
	if !ok {
		return id, "", "", false
	}
	return id, r[0], r[1], true
}

// return id, name, format, ok
func GetServerMeta(id int) (int, string, string, bool) {

	if instance == nil {
		NewMeta()
	}

	r, ok := instance.servers[id]
	if !ok {
		return id, "", "", false
	}
	return id, r[0], r[1], true
}

func (meta *Meta) parseBody(body []byte) {
	/* FORMATS */

	server_formats := meta.getFormats(body, meta.ServerFormats.Start, meta.ServerFormats.End, meta.ServerFormats.Pattern)
	if len(server_formats) == 0 {
		meta.Error = fmt.Errorf("[net.getFormats] getFormats by server return zero len length. start: %#v, end: %#v, pattern: %#v", string(meta.ServerFormats.Start), string(meta.ServerFormats.End), string(meta.ServerFormats.Pattern))
		return
	}

	client_formats := meta.getFormats(body, meta.ClientFormats.Start, meta.ClientFormats.End, meta.ClientFormats.Pattern)
	if len(client_formats) == 0 {
		meta.Error = fmt.Errorf("[net.getFormats] getFormats by client return zero length. start %s, end: %s, pattern: %s", meta.ClientFormats.Start, meta.ClientFormats.End, meta.ClientFormats.Pattern)
		return
	}

	/* TYPES */

	server_types := meta.getTypes(body, meta.ServerTypes.Start, meta.ServerTypes.End, meta.ServerTypes.Pattern)
	if len(server_types) == 0 {
		meta.Error = fmt.Errorf("[net.getTypes] getTypes by server return zero length. start %s, end: %s, pattern: %s", meta.ServerTypes.Start, meta.ServerTypes.End, meta.ServerTypes.Pattern)
		return
	}

	client_types := meta.getTypes(body, meta.ClientTypes.Start, meta.ClientTypes.End, meta.ClientTypes.Pattern)
	if len(client_types) == 0 {
		meta.Error = fmt.Errorf("[net.getTypes] getTypes by client return zero length. start %s, end: %s, pattern: %s", meta.ClientTypes.Start, meta.ClientTypes.End, meta.ClientTypes.Pattern)
		return
	}

	/* GENERATE */

	meta.servers = meta.generateData(server_formats, server_types)
	if len(meta.servers) == 0 {
		meta.Error = fmt.Errorf("[net.generateData] generateData by server return zero length. formats len: %d, types len: %d", len(server_formats), len(server_types))
		return
	}

	meta.clients = meta.generateData(client_formats, client_types)
	if len(meta.clients) == 0 {
		meta.Error = fmt.Errorf("[net.generateData] generateData by client return zero length. formats len: %d, types len: %d", len(client_formats), len(client_types))
		return
	}

	fmt.Printf("meta parse success! servers: %d, clients: %d\n", len(meta.servers), len(meta.clients))
}

func (meta *Meta) generateData(formats []string, types map[int]string) map[int][2]string {
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

func (meta *Meta) getFormats(body []byte, start []byte, end []byte, pattern []byte) []string {

	body = meta.split(body, start, end)
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

func (meta *Meta) getTypes(body []byte, start []byte, end []byte, pattern []byte) map[int]string {
	body = meta.split(body, start, end)
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

func (meta *Meta) getBody() []byte {
	url := fmt.Sprintf("%s%s%s", meta.HostPath, meta.Version, meta.ScriptPath)
	body, err := meta.request(url)
	if err != nil {
		meta.Error = fmt.Errorf("[net.getBody] fail get body. %s", err)
		return nil
	}
	return body
}

func (meta *Meta) getVersion() {

	url := fmt.Sprintf("%s%s", meta.HostPath, meta.VersionPath)

	body, err := meta.request(url)
	if err != nil {
		meta.Error = fmt.Errorf("[net.getVersion] fail get version. %s", err)
		return
	}

	type version struct {
		V string `json:"browser"`
	}

	v := version{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		meta.Error = fmt.Errorf("[net.getVersion] fail get version. %s", err)
		return
	}

	if len(v.V) == 0 {
		meta.Error = fmt.Errorf("[net.getVersion] fail get version. %s", ErrEmptyResult)
		return
	}

	meta.Version = v.V
}

func (meta *Meta) request(url string) ([]byte, error) {

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

func (meta *Meta) split(body []byte, start []byte, end []byte) []byte {
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
