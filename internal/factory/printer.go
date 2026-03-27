package factory

import (
	"fmt"
	"strings"
)

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
	if IsOpen(chest.GetId()) {
		state = "open"
		name = fmt.Sprintf("%s%s%s", green, chest.GetName(), reset)
	}
	fmt.Printf(" %s  %-*s %-*s %-*s %s\n", chest.GetEmoji(), chestNameWidth+len(name)-len(chest.GetName()), name, chestKindWidth, chest.GetKind(), chestStateWidth, state, chest.GetDescription())
}

func GetChestStringForSelect(chest Chest) string {
	return fmt.Sprintf(" %s  %s", chest.GetEmoji(), chest.GetName())
}

func GetJewelStringForSelect(jewel Jewel, chestName string) string {
	return fmt.Sprintf(" %s  %s(%s)", jewel.GetEmoji(), jewel.GetName(), chestName)
}
