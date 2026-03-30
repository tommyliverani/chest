package factory

import (
	"chest/internal/common"
	"encoding/json"
	"fmt"
	"slices"
)

type JewelCreator func(name string, description string) (Jewel, error)
type JewelParser func(data json.RawMessage) (Jewel, error)

var jewelCreatorRegistry = make(map[string]JewelCreator)
var jewelParserRegistry = make(map[string]JewelParser)

func RegisterJewelCreator(kind string, creator JewelCreator) {
	jewelCreatorRegistry[kind] = creator
}

func RegisterJewelParser(kind string, parser JewelParser) {
	jewelParserRegistry[kind] = parser
}

func CreateJewel(kind string, name string, description string) (Jewel, error) {
	creator, exists := jewelCreatorRegistry[kind]
	if !exists {
		return nil, fmt.Errorf("unknown jewel kind: %s", kind)
	}
	return creator(name, description)
}

func ParseJewel(data json.RawMessage) (Jewel, error) {
	kind, err := common.GetKindFromJson(data)
	common.CheckWithMsg("Error while getting jewel kind", err)
	parser, exists := jewelParserRegistry[kind]
	if !exists {
		return nil, fmt.Errorf("unknown jewel kind: %s", kind)
	}
	return parser(data)
}

// a jewel is available if it has both a creator and a parser registered
func GetAvailableJewelKinds() []string {
	availableKinds := make([]string, 0, len(jewelParserRegistry))
	for kind := range jewelParserRegistry {
		if _, exists := jewelCreatorRegistry[kind]; exists {
			availableKinds = append(availableKinds, kind)
		}
	}
	slices.Sort(availableKinds)
	return availableKinds
}

func CreateKeyJewel(chest Chest) (Jewel, error) {
	kind := chest.GetKeyJewelKind()
	return CreateJewel(kind, "keyJewelFor"+chest.GetName(), "keyJewelDescription")
}
