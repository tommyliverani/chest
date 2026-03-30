package factory

import (
	"chest/internal/common"
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
	kind, err := common.GetKindFromJson(data)
	common.CheckWithMsg("Error while getting chest kind", err)
	parser, exists := chestParserRegistry[kind]
	if !exists {
		return nil, fmt.Errorf("unknown chest kind: %s", kind)
	}
	return parser(data)
}

func GetChestById(id string) (Chest, error) {
	chestJson, err := common.GetChestJsonById(id)
	common.CheckWithMsg("Error while retrieving chest JSON", err)
	return ParseChest(chestJson)
}

func SaveOrUpdateChest(chest Chest) (string, error) {
	jsonChest, err := chest.ToJson()
	common.CheckWithMsg("Error while converting chest to JSON", err)
	chestPath, err := common.CreateOrUpdateJsonChestFile(chest.GetId(), jsonChest)
	common.CheckWithMsg("Error while writing chest file", err)
	return chestPath, nil
}

// a chest is available if it has both a creator and a parser registered
func GetAvailableChestKinds() []string {
	availableKinds := make([]string, 0, len(chestParserRegistry))
	for kind := range chestParserRegistry {
		if _, exists := chestCreatorRegistry[kind]; exists {
			availableKinds = append(availableKinds, kind)
		}
	}
	slices.Sort(availableKinds)
	return availableKinds
}

func CheckChestName(oldName string, newName string) error {
	if oldName != newName && slices.Contains(GetExistingChestNames(), newName) {
		return fmt.Errorf("a chest with the name %q already exists", newName)
	}
	return nil
}
