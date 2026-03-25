package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
)

func EditChest() {
	chestNames, err := common.GetExistingChestNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest names: %v\n", err)
		os.Exit(1)
	}
	if len(chestNames) == 0 {
		fmt.Println("No chests available to edit.")
		return
	}
	chestToEdit, err := common.SelectField("Select chest to edit: ", chestNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting chest to edit: %v\n", err)
		os.Exit(1)
	}
	EditChestByName(chestToEdit)
}

func EditChestByName(name string) {
	chest, err := factory.GetExistingChest(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest: %v\n", err)
		os.Exit(1)
	}
	err = chest.Edit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error editing chest: %v\n", err)
		os.Exit(1)
	}
	//TODO: manage chest reaniming if name is changed
	path, err := factory.SaveOrUpdateChest(chest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving chest: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Chest updated in %s\n", path)
}
