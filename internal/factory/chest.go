package factory

import (
	"encoding/json"
)

//ok

type Chest interface {
	GetId() string
	GetName() string
	GetKind() string
	GetDescription() string
	GetEmoji() string
	Delete() error
	Edit(keyJewel Jewel) error
	GetJewels(keyJewel Jewel) ([]Jewel, error)
	GetKeyJewelKind() string // Get the kind of the key jewel required to access the chest
	AddJewel(jewelToAdd Jewel, keyJewel Jewel) error
	UpdateJewel(jewelName string, newJewel Jewel, keyJewel Jewel) error
	RemoveJewel(jewel Jewel, keyJewel Jewel) error
	ToJson() (json.RawMessage, error)
	Open(jewel Jewel) error
	Close() error
}
