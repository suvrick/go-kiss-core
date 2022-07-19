package until

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrEmptyParams = errors.New("empty params")
	ErrEmptyResult = errors.New("empty result")
	ErrBadRequest  = errors.New("bad request")
)

func Request(url string) ([]byte, error) {

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

	return body, nil
}

func RemoveSpacialSymbol(body []byte) []byte {
	body = bytes.ReplaceAll(body, []byte{9}, []byte{})  // remove "	"
	body = bytes.ReplaceAll(body, []byte{32}, []byte{}) // remove " "
	body = bytes.ReplaceAll(body, []byte{10}, []byte{}) // remove "\n"
	body = bytes.ReplaceAll(body, []byte{13}, []byte{}) // remove "\r"
	return body
}
