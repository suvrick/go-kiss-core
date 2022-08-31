package scheme

import (
	"github.com/suvrick/go-kiss-core/meta"
)

type Scheme struct {
	MetaID     uint
	FeildName  string
	FeildIndex int
	FeildType  int
}

func Fill(indexData int, indexFormat int, payload []interface{}, schemes []Scheme, meta *meta.Meta, result map[string]interface{}) error {

	for _, code := range meta.Format {

		indexFormat++

		if code == ',' {
			continue
		}

		if code == '[' {
			Fill(indexData, indexFormat, payload[indexData:], schemes, meta, result)
			continue
		}

		if code == ']' {
			continue
		}

		for _, scheme := range schemes {
			if index != scheme.FeildIndex {
				continue
			}

			result[scheme.FeildName] = payload[scheme.FeildIndex]
		}
	}

	return nil
}

func B(data map[string]interface{}, schemes []Scheme, meta *meta.Meta) ([]interface{}, error) {

	payload := make([]interface{}, 0)

	for i, code := range meta.Format {

		if code == ',' {
			continue
		}

		if code == '[' {
			continue
		}

		if code == ']' {
			continue
		}

		for _, scheme := range schemes {
			if i != scheme.FeildIndex {
				continue
			}

			p := data[scheme.FeildName]

			payload = append(payload, p)
		}

		payload = append(payload, nil)
	}
	return payload, nil
}
