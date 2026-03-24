package command

import (
	. "chest/internal/chest"
	. "chest/internal/common"
	"fmt"
	"os"
)

func EditChest() {
	name, _ := SelectChest(GetChestNames())
	EditChestByName(name)
}

func EditChestByName(name string) {
	chest := GetExistingChestByName(name)
	if chest != nil {
		chest.Edit()
		chestPath := ChestBasePath + "/" + name + ".json"
		err := SaveChestToFile(chestPath, chest)
		if err != nil {
			fmt.Printf("Error saving chest: %v\n", err)
		}
	}
}

func GetExistingChestByName(name string) Chest {
	chestPath := ChestBasePath + "/" + name + ".json"
	chestFile, errRead := os.ReadFile(chestPath)
	if errRead == nil {
		chest, errParse := ParseChest(chestFile)
		if errParse == nil {
			return chest
		} else {
			fmt.Printf("failed to parse chest file %s: %v\n", chestPath, errParse)
		}
	} else {
		fmt.Printf("failed to read chest file %s: %v\n", chestPath, errRead)
	}
	return nil
}
