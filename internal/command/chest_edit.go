package command

import (
	. "chest/internal/chest"
	. "chest/internal/common"
	"fmt"
)

func EditChest() {
	name, _ := SelectChest(GetChestNames())
	EditChestByName(name)
}

func EditChestByName(name string) {
	chest := GetChestByName(name)
	if chest != nil {
		chest.Edit()
		chestPath := ChestBasePath + "/" + name + ".json"
		err := SaveChestToFile(chestPath, chest)
		if err != nil {
			fmt.Printf("Error saving chest: %v\n", err)
		}
	}
}
