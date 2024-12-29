package utils

import (
	"encoding/json"
	"io"
)

func ParseJsonBody(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}

func StrToPointer(s string) *string {
	return &s
}
