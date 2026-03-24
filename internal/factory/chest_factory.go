package factory

import (
	"encoding/json"
	"fmt"
	"slices"
)

type ChestCreator func(name string, description string) (Chest, error)
type ChestParser func(data json.RawMessage) (Chest, error)

var chestCreatorRegistry = make(map[string]ChestCreator)
var chestParserRegistry = make(map[string]ChestParser)

func RegisterChestCreator(kind string, creator ChestCreator) {
	chestCreatorRegistry[kind] = creator
}

func RegisterChestParser(kind string, parser ChestParser) {
	chestParserRegistry[kind] = parser
}

func CreateChest(kind string, name string, description string) (Chest, error) {
	creator, exists := chestCreatorRegistry[kind]
	if !exists {
		return nil, fmt.Errorf("unknown chest kind: %s", kind)
	}
	return creator(name, description)
}

func ParseChest(data json.RawMessage) (Chest, error) {
	var helper struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(data, &helper); err != nil {
		return nil, fmt.Errorf("error while parsing chest: %w", err)
	}
	parser, exists := chestParserRegistry[helper.Kind]
	if !exists {
		return nil, fmt.Errorf("unknown chest kind: %s", helper.Kind)
	}
	return parser(data)
}

// a chest is available if it has both a creator and a parser registered
func GetChestKinds() []string {
	availableKinds := make([]string, 0, len(chestParserRegistry))
	for kind := range chestParserRegistry {
		if _, exists := chestCreatorRegistry[kind]; exists {
			availableKinds = append(availableKinds, kind)
		}
	}
	slices.Sort(availableKinds)
	return availableKinds
}
