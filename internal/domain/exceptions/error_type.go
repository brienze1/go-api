package exceptions

import (
	"encoding/json"
	"strings"
)

type ErrorType struct {
	StatusCode int      `json:"status_code"`
	Messages   []string `json:"messages"`
	Type       string   `json:"type"`
}

func (e ErrorType) Error() string {
	return strings.Join(e.Messages, "; ")
}

func (e ErrorType) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}
