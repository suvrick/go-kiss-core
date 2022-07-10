package meta

import "errors"

var (
	ErrEmptyParams  = errors.New("empty params")
	ErrEmptyResult  = errors.New("empty result")
	ErrGenerateMeta = errors.New("fail generate meta")
	ErrBadRequest   = errors.New("bad request")
)
