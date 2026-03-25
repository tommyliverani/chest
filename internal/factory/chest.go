package factory

import (
	"encoding/json"
)

type Chest interface {
	GetName() string
	GetKind() string
	GetDescription() string
	GetEmoji() string
	Delete() error
	Edit() error
	GetJewels(keyJewel json.RawMessage) ([]Jewel, error)
	AddJewel(jewelToAdd Jewel, keyJewel json.RawMessage) error
	RemoveJewel(jewelName string, keyJewel json.RawMessage) error
	ToJson() (json.RawMessage, error)
	Open() (Jewel, error)
	Close() error
}
