package command

import (
	. "chest/internal/chest"
	. "chest/internal/common"
	"fmt"
)

func DeleteChest() {
	chestToDelete, _ := SelectChest(GetChestNames())
	DeleteChestByName(chestToDelete)
}

func DeleteChestByName(chestToDelete string) {
	chest := GetChestByName(chestToDelete)
	if chest != nil {
		err := chest.Delete()
		if err == nil {
			fmt.Printf("Chest %s deleted\n", chestToDelete)
		}
	}
}
