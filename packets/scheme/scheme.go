package scheme

import "github.com/suvrick/go-kiss-core/until"

type Scheme struct {
	PacketID int         `json:"packet_id"`
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	Type     int         `json:"type"`
	Index    int         `json:"index"`
}

type Schemes struct {
	schemes []Scheme
	Error   error
}

var Instance *Schemes

func NewSchemes(path string) *Schemes {

	s := &Schemes{}
	if err := until.LoadConfigFromFile(path, s); err != nil {
		s.Error = err
		return s
	}

	return s
}

func NewDefaultSchemes() *Schemes {
	return &Schemes{
		schemes: []Scheme{
			{
				PacketID: 4,
				Name:     "result",
				Index:    0,
			},
			{
				PacketID: 4,
				Name:     "user_id",
				Index:    1,
			},
			{
				PacketID: 4,
				Name:     "balance",
				Index:    2,
			},
		},
	}
}

func (s *Schemes) Fill(packetID int, source map[string]interface{}, params []interface{}) map[string]interface{} {
	for _, v := range s.FindAllSchemesByPacketID(packetID) {
		if v.Index < len(params) {
			source[v.Name] = params[v.Index]
		}
	}
	return source
}

func (s *Schemes) Load(packetID int, source map[string]interface{}, params []interface{}) []interface{} {
	for _, v := range s.FindAllSchemesByPacketID(packetID) {
		value, ok := source[v.Name]
		if ok {
			if v.Index < len(params) {
				params[v.Index] = value
			}
		}
	}
	return params
}

func (s *Schemes) Add(scheme Scheme) {
	for i, v := range s.schemes {
		if scheme.Equle(v) {
			s.schemes[i] = scheme
			return
		}
	}

	s.schemes = append(s.schemes, scheme)
}

func (s *Schemes) Remove(scheme Scheme) int {
	for i, v := range s.schemes {
		if scheme.Equle(v) {
			s.schemes[i] = Scheme{}
		}
	}
	return -1
}

func (s *Schemes) FindAllSchemesByPacketID(packetID int) []Scheme {
	result := make([]Scheme, 0)
	for _, v := range s.schemes {
		if v.PacketID == packetID {
			result = append(result, v)
		}
	}
	return result
}

func (s *Schemes) Find(scheme Scheme) []Scheme {
	result := make([]Scheme, 0)
	for _, v := range s.schemes {
		if scheme.Equle(v) {
			result = append(result, v)
		}
	}
	return result
}

func (s *Schemes) IsContaines(scheme Scheme) bool {
	return len(s.Find(scheme)) == 0
}

func (s *Scheme) Equle(s2 Scheme) bool {
	return s.Index == s2.Index && s.Name == s2.Name && s.PacketID == s2.PacketID
}
