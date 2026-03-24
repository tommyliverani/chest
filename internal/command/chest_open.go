package command

import (
	. "chest/internal/common"
	. "chest/internal/jewel"
	. "chest/internal/session"
	"fmt"
)

func OpenChestByName(chestName string) {
	chest := GetExistingChestByName(chestName)
	jewelKey, err := CreateJewel(chest.GetKeyJewelKind(), chestName+"JewelKey", "This is the key jewel for "+chestName)
	if err != nil {
		fmt.Printf("Error creating key jewel: %v\n", err)
		return
	}
	jewelData, err := jewelKey.ToJson()
	if err != nil {
		fmt.Printf("Error converting key jewel to JSON: %v\n", err)
		return
	}
	jewels, err := chest.GetJewels(jewelData)
	if err != nil {
		fmt.Printf("Error opening chest: %v\nerror:%v\n", chestName, err)
		return
	}
	StoreSession(chestName, jewelData)
	PrintJewels(jewels)
}

func OpenChest() {
	chestToOpen, _ := SelectChest(GetChestNames())
	OpenChestByName(chestToOpen)
}
