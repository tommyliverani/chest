package factory

import (
	"chest/internal/common"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
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

func GetExistingChest(name string) (Chest, error) {

	chestJson, err := common.GetExistingChestJson(name)
	if err != nil {
		return nil, err
	}
	return ParseChest(chestJson)
}

func SaveOrUpdateChest(chest Chest) (string, error) {
	jsonChest, err := chest.ToJson()
	if err != nil {
		return "", fmt.Errorf("error converting chest to JSON: %w", err)
	}
	chestPath, err := common.CreateOrUpdateJsonChestFile(chest.GetName(), jsonChest)
	if err != nil {
		return "", fmt.Errorf("error saving chest: %w", err)
	}
	return chestPath, nil
}

const (
	chestNameWidth  = 16
	chestKindWidth  = 12
	chestStateWidth = 10
	chestLineWidth  = 3 + chestNameWidth + 1 + chestKindWidth + 1 + chestStateWidth + 1 + len("DESCRIPTION")
)

func PrintChestHeader() {
	fmt.Printf("     %-*s %-*s %-*s %s\n", chestNameWidth, "NAME", chestKindWidth, "KIND", chestStateWidth, "STATE", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", chestLineWidth))
}

func PrintChestFooter(count int) {
	fmt.Println(strings.Repeat("-", chestLineWidth))
	fmt.Printf("%d chest(s) available\n", count)
}

const green = "\033[32m"
const reset = "\033[0m"

func PrintChest(chest Chest) {
	state := "close"
	name := chest.GetName()
	if IsOpen(chest.GetName()) {
		state = "open"
		name = fmt.Sprintf("%s%s%s", green, chest.GetName(), reset)
	}
	fmt.Printf(" %s  %-*s %-*s %-*s %s\n", chest.GetEmoji(), chestNameWidth+len(name)-len(chest.GetName()), name, chestKindWidth, chest.GetKind(), chestStateWidth, state, chest.GetDescription())
}

func GetChestString(chest Chest) string {
	state := "close"
	name := chest.GetName()
	if IsOpen(chest.GetName()) {
		state = "open"
	}
	return fmt.Sprintf(" %s  %-*s %-*s %-*s %s\n", chest.GetEmoji(), chestNameWidth+len(name)-len(chest.GetName()), name, chestKindWidth, chest.GetKind(), chestStateWidth, state, chest.GetDescription())
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
