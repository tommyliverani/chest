package chest

import (
	"encoding/json"
	"fmt"
	"slices"
)

type ChestCreator func(name string, description string) (Chest, error)

var creatorRegistry map[string]ChestCreator

func ChestCreatorRegistry() map[string]ChestCreator {
	if creatorRegistry == nil {
		creatorRegistry = make(map[string]ChestCreator)
	}
	return creatorRegistry
}

func RegisterChestCreator(kind string, creator ChestCreator) {
	ChestCreatorRegistry()[kind] = creator
}

func GetChestCreator(kind string) (ChestCreator, bool) {
	creator, exists := ChestCreatorRegistry()[kind]
	return creator, exists
}

type ChestParser func(data json.RawMessage) (Chest, error)

var parserRegistry map[string]ChestParser

func ChestParserRegistry() map[string]ChestParser {
	if parserRegistry == nil {
		parserRegistry = make(map[string]ChestParser)
	}
	return parserRegistry
}

func RegisterChestParser(kind string, parser ChestParser) {
	ChestParserRegistry()[kind] = parser
}

func GetChestParser(kind string) (ChestParser, bool) {
	parser, exists := ChestParserRegistry()[kind]
	return parser, exists
}

func GetKinds() []string {
	availableKinds := make([]string, 0, len(parserRegistry))
	for kind := range parserRegistry {
		if _, exists := creatorRegistry[kind]; exists {
			availableKinds = append(availableKinds, kind)
		}
	}
	slices.Sort(availableKinds)
	return availableKinds
}

func ParseChest(data json.RawMessage) (Chest, error) {
	var helper struct {
		Kind string `json:"kind"`
	}
	json.Unmarshal(data, &helper)

	parser, exists := GetChestParser(helper.Kind)
	if !exists {
		return nil, fmt.Errorf("unknown chest kind: %s", helper.Kind)
	}
	return parser(data)
}
