package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
	"slices"
)

func OpenChestByName(chestName string) {

	openedChests, err := factory.GetOpenChestNames()
	if err != nil {
		fmt.Printf("Error retrieving session: %v\n", err)
	}

	if openedChests != nil && slices.Contains(openedChests, chestName) {
		fmt.Fprintf(os.Stderr, "Chest %s is already open\n", chestName)
		os.Exit(1)
	}
	chest, err := factory.GetExistingChest(chestName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest: %v\n", err)
		os.Exit(1)
	}
	if chest == nil {
		fmt.Fprintf(os.Stderr, "Chest '%s' not found\n", chestName)
		os.Exit(1)
	}
	jewel, err := chest.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening chest: %v\n", err)
		os.Exit(1)
	}
	jewelData, err := jewel.ToJson()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting jewel to JSON: %v\n", err)
		os.Exit(1)
	}
	err = factory.StoreSession(chestName, jewelData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error storing session: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Chest %s opened successfully\n", chestName)
	//todo: aggiungere la print dei jewel?
}

func OpenChest() {
	chestNames, err := common.GetExistingChestNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest names: %v\n", err)
		os.Exit(1)
	}
	openedChests, err := factory.GetOpenChestNames()
	if err != nil {
		fmt.Printf("Error retrieving session: %v\n", err)
	}
	if openedChests != nil {
		filtered := make([]string, 0, len(chestNames))
		for _, name := range chestNames {
			if !slices.Contains(openedChests, name) {
				filtered = append(filtered, name)
			}
		}
		chestNames = filtered
	}
	if len(chestNames) == 0 {
		fmt.Println("No chest available to open")
		return
	}

	chestToOpen, _ := common.SelectField("Select chest to open: ", chestNames)
	OpenChestByName(chestToOpen)
}
