package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
)

func CloseChest() {

	existingChests := factory.GetAllOpenChests()
	if len(existingChests) == 0 {
		fmt.Println("No chests available to close")
		return
	}
	chestToClose := factory.SelectChest("Select chest to close: ", existingChests)
	closeChest(chestToClose)
}

func CloseChestByName(name string) {
	chestToClose, found := factory.FindChestByName(name)
	if !found {
		common.PrintErrorAndExit(fmt.Sprintf("Chest '%s' not found", name))
	}
	closeChest(chestToClose)
}

func closeChest(chestToClose factory.Chest) {
	err := chestToClose.Close()
	common.Check(err)
	factory.DeleteSession(chestToClose.GetId())
	fmt.Printf("%s closed\n", chestToClose.GetName())
}
