package net

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
		ServerFormats: RegConfig{
			Start:   []byte("PacketServer=o,o.FORMATS=["),
			End:     []byte("}"),
			Pattern: []byte(`"([A-Z,\[\]]*)"`),
		},
		ClientFormats: RegConfig{
			Start:   []byte("PacketServer=o,o.FORMATS=["),
			End:     []byte("}"),
			Pattern: []byte(`"([A-Z,\[\]]*)"`),
		},
		ServerTypes: RegConfig{
			Start:   []byte(".ServerPacketType=void"),
			End:     []byte("}"),
			Pattern: []byte(`.([A-Z_]*)=([0-9]*)`),
		},
		ClientTypes: RegConfig{
			Start:   []byte(".ClientPacketType=void"),
			End:     []byte("}"),
			Pattern: []byte(`.([A-Z_]*)=([0-9]*)`),
		},
	}
}
