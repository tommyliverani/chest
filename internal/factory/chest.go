package factory

import (
	. "chest/internal/jewel"
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
	AddJewel(jewel Jewel) error
	RemoveJewel(jewelName string) error
	ToJson() (json.RawMessage, error)
	GetKeyJewelKind() string
}
