package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
)

func OpenChestByName(chestName string) {
	chest, found := factory.FindChestByName(chestName)
	if !found {
		common.PrintErrorAndExit(fmt.Sprintf("Chest '%s' not found", chestName))
	}
	openChest(chest)
}

func AskNameAndOpenChest() {
	existingChests := factory.GetAllChests()
	if len(existingChests) == 0 {
		fmt.Println("No chests available to open.")
		return
	}
	chestToOpen := factory.SelectChest("Select chest to open: ", existingChests)
	openChest(chestToOpen)
}

func openChest(chest factory.Chest) {
	keyJewel, err := factory.CreateKeyJewel(chest)
	common.Check(err)
	err = chest.Open(keyJewel)
	common.Check(err)
	factory.StoreSession(chest.GetId(), keyJewel)
	fmt.Printf("%s opened\n", chest.GetName())
}
