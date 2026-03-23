package command

import (
	. "chest/internal/chest"
)

func ListChests() {
	PrintChestHeader()
	for _, chest := range GetAllChests() {
		PrintChest(chest)
	}
}
