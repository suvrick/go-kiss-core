package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/suvrick/go-kiss-core/until"
)

var Instance *Meta

type Meta struct {
	HostPath    string `json:"host_path"`
	VersionPath string `json:"version_path"`
	ScriptPath  string `json:"script_path"`
	Version     string `json:"version"`

	ClientFormats RegConfig `json:"client_formats"`
	ServerFormats RegConfig `json:"server_formats"`
	ServerTypes   RegConfig `json:"server_types"`
	ClientTypes   RegConfig `json:"client_types"`

	Clients map[uint16][2]string `json:"clients"`
	Servers map[uint16][2]string `json:"servers"`

	HasChangedVersion bool  `json:"-"`
	Error             error `json:"-"`
}

type RegConfig struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Pattern string `json:"pattern"`
}

func NewMeta(path string) *Meta {

	meta := &Meta{}
	if err := until.LoadConfigFromFile(path, meta); err != nil {
		meta.Error = err
		return meta
	}

	if meta.CheckUpdateMeta() {
		if err := until.SaveConfigToFile(path, meta); err != nil {
			meta.Error = err
			return meta
		}
	}

	return meta
}

func (meta *Meta) CheckUpdateMeta() bool {

	log.Printf("check new version...\n")

	v := meta.getVersion()
	if len(v) == 0 {
		meta.Error = fmt.Errorf("error check new version")
		return false
	}

	if v == meta.Version {
		log.Printf("install actual version %s\n", v)
		return false
	}

	log.Printf("get new version %s\n", v)

	meta.Version = v

	log.Printf("try install new version\n")

	meta.Initialize()

	if meta.Error != nil {
		log.Printf("new version install FAIL\n")
		return false
	}

	log.Printf("new version install OK\n")
	return true
}

func (meta *Meta) Initialize() {

	body := meta.getBody()
	if meta.Error != nil {
		return
	}

	meta.parseBody(body)
}

// return name, format, ok
func (meta *Meta) GetClientMeta(typeID uint16) (string, string, bool) {

	if meta.Clients == nil {
		return "", "", false
	}

	r, ok := meta.Clients[typeID]
	if !ok {
		return "", "", false
	}
	return r[0], r[1], true
}

// return name, format, ok
func (meta *Meta) GetServerMeta(typeID uint16) (string, string, bool) {

	if meta.Servers == nil {
		return "", "", false
	}

	r, ok := meta.Servers[typeID]
	if !ok {
		return "", "", false
	}
	return r[0], r[1], true
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

	meta.Servers = meta.generateData(server_formats, server_types)
	if len(meta.Servers) == 0 {
		meta.Error = fmt.Errorf("[net.generateData] generateData by server return zero length. formats len: %d, types len: %d", len(server_formats), len(server_types))
		return
	}

	meta.Clients = meta.generateData(client_formats, client_types)
	if len(meta.Clients) == 0 {
		meta.Error = fmt.Errorf("[net.generateData] generateData by client return zero length. formats len: %d, types len: %d", len(client_formats), len(client_types))
		return
	}
}

func (meta *Meta) generateData(formats []string, types map[uint16]string) map[uint16][2]string {
	result := make(map[uint16][2]string)
	for id, format := range formats {
		name, ok := types[uint16(id)]
		if ok {
			result[uint16(id)] = [2]string{
				name,
				format,
			}
		}
	}

	return result
}

func (meta *Meta) getFormats(body []byte, start string, end string, pattern string) []string {

	body = meta.split(body, []byte(start), []byte(end))
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

func (meta *Meta) getTypes(body []byte, start string, end string, pattern string) map[uint16]string {
	body = meta.split(body, []byte(start), []byte(end))
	if len(body) == 0 {
		return nil
	}

	//fmt.Printf("body types: %s", body)

	reg, err := regexp.Compile(string(pattern))
	if err != nil {
		return nil
	}

	finded := reg.FindAllSubmatch(body, -1)

	result := make(map[uint16]string, len(finded))

	for _, v := range finded {

		if len(v) != 3 {
			continue
		}

		key := v[2]
		value := v[1]
		id, err := strconv.Atoi(string(key))
		if err == nil {
			result[uint16(id)] = string(value)
		}
	}

	return result
}

func (meta *Meta) getBody() []byte {
	url := fmt.Sprintf("%s%s%s", meta.HostPath, meta.Version, meta.ScriptPath)
	body, err := until.Request(url)
	if err != nil {
		meta.Error = fmt.Errorf("[net.getBody] fail get body. %s", err)
		return nil
	}

	return until.RemoveSpacialSymbol(body)
}

func (meta *Meta) getVersion() string {

	url := fmt.Sprintf("%s%s", meta.HostPath, meta.VersionPath)

	body, err := until.Request(url)

	if err != nil {
		return ""
	}

	type version struct {
		V string `json:"browser"`
	}

	v := version{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return ""
	}

	if len(v.V) == 0 {
		return ""
	}

	return v.V
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
