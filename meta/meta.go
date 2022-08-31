package meta

import "io"

type MetaType byte

const (
	CLIENT MetaType = 0
	SERVER MetaType = 1
)

type Meta struct {
	ID     uint
	Format string
	Name   string
	Type   MetaType
}

func Parser(reader io.Reader, meta *Meta) ([]interface{}, error) {
	return nil, nil
}

func Fill(payload []interface{}, meta *Meta) ([]byte, error) {
	return nil, nil
}
