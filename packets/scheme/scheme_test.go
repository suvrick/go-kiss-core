package scheme

import (
	"testing"
)

func TestNewSchemes(t *testing.T) {
	schemes := Schemes{
		schemes: []Scheme{
			{
				PacketID: 4,
				Name:     "result",
				Index:    0,
			},
			{

				PacketID: 4,
				Name:     "userID",
				Index:    5,
			},
		},
	}

	source := make(map[string]interface{}, 0)

	source = schemes.Fill(4, source, []interface{}{0, 23, 432, "asd"})

}
