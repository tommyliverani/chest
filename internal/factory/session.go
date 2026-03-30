package factory

import (
	"encoding/json"
)

type Session struct {
	KeyJewel json.RawMessage `json:"jewel"`
}

type SessionMap = map[string]Session // chestId -> Session
