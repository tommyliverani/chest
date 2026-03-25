package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
)

func AddJewel(kind string) {
	name, err := common.ReadField("Insert jewel name: ")
	if err != nil {
		fmt.Printf("Error reading jewel name: %v\n", err)
		return
	}
	AddJewelByName(kind, name)
}

func AddJewelByName(kind string, name string) {
	openChest, err := factory.GetOpenChestNames()
	if err != nil {
		fmt.Printf("Error retrieving open chests: %v\n", err)
		return
	}
	if len(openChest) == 0 {
		fmt.Println("No open chests available. Please open a chest before adding a jewel.")
		return
	}
	var selectedChest string
	if len(openChest) > 1 {
		selectedChest, err = common.SelectField("Select chest to add the jewel to: ", openChest)
		if err != nil {
			fmt.Printf("Error selecting chest: %v\n", err)
			return
		}
	} else {
		//inserire promt di conferma per aggiungere alla chest aperta
		response, err := common.SelectField(fmt.Sprintf("The jewel will be addet to chest %s, ok? (y/n): ", openChest[0]), []string{"y", "n"})
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}
		if response != "y" && response != "Y" {
			fmt.Println("Aborting jewel addition.")
			return
		}
		selectedChest = openChest[0]
	}
	keyJewel, exists, err := factory.GetKeyJewel(selectedChest)
	if err != nil {
		fmt.Printf("Error retrieving key jewel for chest: %v\n", err)
		return
	}
	if !exists {
		fmt.Println("No key jewel found for the selected chest.")
		return
	}
	description, err := common.ReadField("Insert jewel description: ")
	if err != nil {
		fmt.Printf("Error reading jewel description: %v\n", err)
		return
	}
	chest, err := factory.GetExistingChest(selectedChest)
	if err != nil {
		fmt.Printf("Error retrieving chest: %v\n", err)
		return
	}
	jewel, err := factory.CreateJewel(kind, name, description)
	if err != nil {
		fmt.Printf("Error creating jewel of kind %s: %v\n", kind, err)
		return
	}
	err = chest.AddJewel(jewel, keyJewel)
	if err != nil {
		fmt.Printf("Error adding jewel to chest: %v\n", err)
		return
	}
	_, err = factory.SaveOrUpdateChest(chest)
	if err != nil {
		fmt.Printf("Error saving chest: %v\n", err)
		return
	}
	fmt.Printf("Jewel %s added to chest %s\n", name, selectedChest)
}
