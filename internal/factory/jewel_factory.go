package factory

import (
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

func CreateJewel(kind string, name string, description string) (Jewel, error) {
	creator, exists := jewelCreatorRegistry[kind]
	if !exists {
		return nil, fmt.Errorf("unknown jewel kind: %s", kind)
	}
	return creator(name, description)
}

func RegisterJewelParser(kind string, parser JewelParser) {
	jewelParserRegistry[kind] = parser
}

func ParseJewel(data json.RawMessage) (Jewel, error) {
	var helper struct {
		Kind string `json:"kind"`
	}
	parseError := json.Unmarshal(data, &helper)
	if parseError != nil {
		return nil, fmt.Errorf("error while parsing jewel: %w", parseError)
	}
	parser, exists := jewelParserRegistry[helper.Kind]
	if !exists {
		return nil, fmt.Errorf("unknown jewel kind: %s", helper.Kind)
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
