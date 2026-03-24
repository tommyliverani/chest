package command

import (
	. "chest/internal/chest"
	. "chest/internal/common"
	"fmt"
	"os"
)

func CreateChest() {
	name, _ := ReadChestName()
	CreateChestByName(name)
}

func CreateChestByName(name string) {
	kind, _ := SelectChestKind(GetChestKinds())
	description, _ := ReadChestDescription()
	creator, _ := GetChestCreator(kind)
	chest, _ := creator(name, description)
	chestPath := ChestBasePath + "/" + name + ".json"
	err := SaveChestToFile(chestPath, chest)
	if err != nil {
		fmt.Printf("Error creating chest: %v\n", err)
		return
	}
	fmt.Printf("Chest created in %s\n", chestPath)
}

func SaveChestToFile(filename string, chest Chest) error {
	data, err := chest.ToJson()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
