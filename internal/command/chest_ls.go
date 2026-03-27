package command

import (
	"chest/internal/factory"
	"fmt"
)

func ListChests() {
	chests := factory.GetAllChests()
	if len(chests) == 0 {
		fmt.Println("No chest available")
		return
	} else {
		factory.PrintChestHeader()
		for _, chest := range chests {
			factory.PrintChest(chest)
		}
		factory.PrintChestFooter(len(chests))
	}
}
