package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
)

func DeleteChest() {
	chestNames, err := common.GetExistingChestNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest names: %v\n", err)
		os.Exit(1)
	}
	if len(chestNames) == 0 {
		fmt.Println("No chests available to delete.")
		return
	}
	chestToDelete, err := common.SelectField("Select chest to delete: ", chestNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting chest to delete: %v\n", err)
		os.Exit(1)
	}
	DeleteChestByName(chestToDelete)
}

func DeleteChestByName(chestToDelete string) {
	chest, err := factory.GetExistingChest(chestToDelete)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest: %v\n", err)
		os.Exit(1)
	}
	if chest == nil {
		fmt.Fprintf(os.Stderr, "Chest '%s' not found\n", chestToDelete)
		os.Exit(1)
	}
	err = chest.Delete()
	if err == nil {
		fmt.Printf("Chest %s deleted\n", chestToDelete)
	} else {
		fmt.Fprintf(os.Stderr, "Error deleting chest: %v\n", err)
		os.Exit(1)
	}
}
