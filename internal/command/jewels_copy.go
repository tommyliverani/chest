package command

import (
	"chest/internal/factory"
	"fmt"
)

func AskNameAndCopyJewel(jewelKind string) {
	chestJewels := factory.GetAllJewelsByKindFromOpenChests(jewelKind)
	if len(chestJewels) == 0 {
		fmt.Printf("No jewels of kind %s found in open chests\n", jewelKind)
		return
	}
	chestName, jewelToCopy := factory.SelectJewel("Select jewel to copy: ", chestJewels)
	jewelToCopy.Copy()
	fmt.Printf("Jewel '%s' from chest '%s' copied to clipboard\n", jewelToCopy.GetName(), chestName)
}

func CopyJewelByName(jewelKind string, jewelName string) {
	chestJewels := factory.FindJewelsByKindAndNameFromOpenChests(jewelKind, jewelName)
	if len(chestJewels) == 0 {
		fmt.Printf("No jewels of kind %s found in open chests\n", jewelKind)
		return
	}
	var chestName string
	var jewelToCopy factory.Jewel
	if len(chestJewels) > 1 {
		chestName, jewelToCopy = factory.SelectJewel("Select jewel to copy: ", chestJewels)
	} else {
		chestName = chestJewels[0].ChestName
		jewelToCopy = chestJewels[0].Jewel
	}
	jewelToCopy.Copy()
	fmt.Printf("Jewel '%s' from chest '%s' copied to clipboard\n", jewelToCopy.GetName(), chestName)
}
