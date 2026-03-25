package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
	"slices"
)

func CloseChest() {
	chestNames, err := factory.GetOpenChestNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest names: %v\n", err)
		os.Exit(1)
	}
	name, err := common.SelectField("Which chest do you want to close? ", chestNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading chest name: %v\n", err)
		os.Exit(1)
	}
	CloseChestByName(name)
}

func CloseChestByName(name string) {
	OpenChests, err := factory.GetOpenChestNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving open chests: %v\n", err)
		os.Exit(1)
	}
	if OpenChests != nil && !slices.Contains(OpenChests, name) {
		fmt.Fprintf(os.Stderr, "Chest %s is not open\n", name)
		os.Exit(1)
	}
	chest, err := factory.GetExistingChest(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest: %v\n", err)
		os.Exit(1)
	}
	err = chest.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error closing chest: %v\n", err)
		os.Exit(1)
	}
	factory.DeleteSession(name)
	fmt.Printf("Chest %s closed\n", name)
}
