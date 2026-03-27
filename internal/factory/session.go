package factory

import (
	"encoding/json"
)

// ok
type Session struct {
	KeyJewel json.RawMessage `json:"jewel"`
}

type SessionMap = map[string]Session //chestId -> Session
