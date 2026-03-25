package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
)

func CreateChest() {
	name, err := common.ReadField("Insert chest name: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading chest name: %v\n", err)
		os.Exit(1)
	}
	CreateChestByName(name)
}

func CreateChestByName(name string) {
	kind, err := common.SelectField("Select chest kind: ", factory.GetAvailableChestKinds())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting chest kind: %v\n", err)
		os.Exit(1)
	}
	description, err := common.ReadField("Insert chest description: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading chest description: %v\n", err)
		os.Exit(1)
	}
	chest, err := factory.CreateChest(kind, name, description)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating chest of kind %s: %v\n", kind, err)
		os.Exit(1)
	}
	chestPath, err := factory.SaveOrUpdateChest(chest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating chest: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Chest created in %s\n", chestPath)
}
