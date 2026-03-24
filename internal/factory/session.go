package factory

import (
	"encoding/json"
)

type Session struct {
	Jewel json.RawMessage `json:"jewel"`
}

type SessionMap = map[string]Session
