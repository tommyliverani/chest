package factory

import (
	"encoding/json"
)

//ok

type Jewel interface {
	GetName() string
	GetKind() string
	GetEmoji() string
	GetDescription() string
	ToJson() (json.RawMessage, error)
	Edit() error
	Print()
	Copy()
	Use()
}
