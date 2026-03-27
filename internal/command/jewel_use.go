package command

import (
	"chest/internal/factory"
	"fmt"
)

func UseJewelByName(jewelKind string, jewelName string) {
	chestJewels := factory.FindJewelsByKindAndNameFromOpenChests(jewelKind, jewelName)
	if len(chestJewels) == 0 {
		fmt.Printf("No jewels of kind %s found in open chests\n", jewelKind)
		return
	}
	var jewelToUse factory.Jewel
	if len(chestJewels) > 1 {
		_, jewelToUse = factory.SelectJewel("Select jewel to use: ", chestJewels)
	} else {
		jewelToUse = chestJewels[0].Jewel
	}
	jewelToUse.Use()
}

func AskNameAndUseJewel(jewelKind string) {
	chestJewels := factory.GetAllJewelsByKindFromOpenChests(jewelKind)
	if len(chestJewels) == 0 {
		fmt.Printf("No jewels of kind %s found in open chests\n", jewelKind)
		return
	}
	_, jewelToUse := factory.SelectJewel("Select jewel to use: ", chestJewels)
	jewelToUse.Use()
}
