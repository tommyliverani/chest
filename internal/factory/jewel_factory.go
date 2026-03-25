package factory

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
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

const (
	jewelChestWidth = 16
	jewelNameWidth  = 16
	jewelKindWidth  = 12
	jewelLineWidth  = 3 + jewelChestWidth + 1 + jewelNameWidth + 1 + jewelKindWidth + 1 + len("DESCRIPTION")
)

func PrintJewelHeader() {
	fmt.Printf("     %-*s %-*s %-*s %s\n", jewelNameWidth, "NAME", jewelChestWidth, "CHEST", jewelKindWidth, "KIND", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", jewelLineWidth))
}

func PrintJewel(jewel Jewel, chestName string) {
	fmt.Printf(" %s  %-*s %-*s %-*s %s\n", jewel.GetEmoji(), jewelNameWidth, jewel.GetName(), jewelChestWidth, chestName, jewelKindWidth, jewel.GetKind(), jewel.GetDescription())
}

func PrintJewelFooter(count int) {
	fmt.Println(strings.Repeat("-", jewelLineWidth))
	fmt.Printf("%d jewel(s) available\n", count)
}

func GetJewelString(jewel Jewel, chestName string) string {
	return fmt.Sprintf(" %s  %-*s %-*s %-*s %s\n", jewel.GetEmoji(), jewelNameWidth, jewel.GetName(), jewelChestWidth, chestName, jewelKindWidth, jewel.GetKind(), jewel.GetDescription())
}

func PrintJewels(jewels []Jewel, chestName string) {
	PrintJewelHeader()
	for _, j := range jewels {
		PrintJewel(j, chestName)
	}
	PrintJewelFooter(len(jewels))
}
