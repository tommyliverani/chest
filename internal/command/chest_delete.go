package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"slices"
)

func DeleteChest() {
	existingChests := factory.GetAllChests()
	if len(existingChests) == 0 {
		fmt.Println("No chests available to delete.")
		return
	}
	chestToDelete := factory.SelectChest("Select chest to delete: ", existingChests)
	deleteChestJsonAndSession(chestToDelete)
}

func DeleteChestByName(chestToDelete string) {
	chestNames := factory.GetExistingChestNames()
	if chestNames != nil && !slices.Contains(chestNames, chestToDelete) {
		common.PrintErrorAndExit(fmt.Sprintf("Chest '%s' not found", chestToDelete))
	}
	chest, found := factory.FindChestByName(chestToDelete)
	if !found {
		common.PrintErrorAndExit(fmt.Sprintf("Chest '%s' not found", chestToDelete))
	}
	deleteChestJsonAndSession(chest)
}

func deleteChestJsonAndSession(chestToDelete factory.Chest) {
	confirm := common.SelectField(fmt.Sprintf("Are you sure you want to delete '%s'?", chestToDelete.GetName()), []string{"No", "Yes"})
	if confirm != "Yes" {
		fmt.Println("Delete cancelled")
		return
	}
	chestToDelete.Delete()
	err := common.DeleteChestJsonById(chestToDelete.GetId())
	common.Check(err)
	factory.DeleteSession(chestToDelete.GetId())
	fmt.Printf("%s deleted\n", chestToDelete.GetName())
}
