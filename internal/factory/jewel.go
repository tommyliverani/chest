package factory

import (
	"encoding/json"
)

type Jewel interface {
	GetName() string
	GetKind() string
	GetEmoji() string
	GetDescription() string
	ToJson() (json.RawMessage, error)
}
