package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"slices"
)

func EditChest() {
	existingChests := factory.GetAllChests()
	if len(existingChests) == 0 {
		fmt.Println("No chests available to edit.")
		return
	}
	chestToEdit := factory.SelectChest("Select chest to edit: ", existingChests)
	editChestJson(chestToEdit)
}

func EditChestByName(chestToEdit string) {
	chestNames := factory.GetExistingChestNames()
	if chestNames != nil && !slices.Contains(chestNames, chestToEdit) {
		common.PrintErrorAndExit(fmt.Sprintf("Chest '%s' not found", chestToEdit))
	}
	chest, found := factory.FindChestByName(chestToEdit)
	if !found {
		common.PrintErrorAndExit(fmt.Sprintf("Chest '%s' not found", chestToEdit))
	}
	editChestJson(chest)
}

func editChestJson(chestToEdit factory.Chest) {
	keyJewel, err := factory.CreateKeyJewel(chestToEdit)
	oldName := chestToEdit.GetName()
	common.Check(err)
	common.Check(chestToEdit.Edit(keyJewel))
	common.Check(factory.CheckChestName(oldName, chestToEdit.GetName()))
	_, err = factory.SaveOrUpdateChest(chestToEdit)
	common.Check(err)
	fmt.Printf("%s edited\n", chestToEdit.GetName())
}
