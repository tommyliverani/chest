package command

import (
	. "chest/internal/chest"
	. "chest/internal/common"
	"fmt"
	"os"
)

func ListChests() {
	PrintChestHeader()
	for _, chest := range GetAllChests() {
		PrintChest(chest)
	}
}

// da rivedere con go routines e jewel
func GetAllChests() []Chest {
	var chests []Chest
	for _, chestPath := range GetChestFilePaths() {
		chestFile, errRead := os.ReadFile(chestPath)
		if errRead == nil {
			chest, errParse := ParseChest(chestFile)
			if errParse == nil {
				chests = append(chests, chest)
			}

		}
	}
	return chests
}

func PrintJewelHeader() {
	fmt.Println("\tNAME\t\tKIND\t\tDESCRIPTION")
}
func PrintJewel(jewel Jewel) {
	fmt.Printf("%s\t%s\t\t%s\t\t%s\n", jewel.GetEmoji(), jewel.GetName(), jewel.GetKind(), jewel.GetDescription())
}

func PrintJewels(jewels []Jewel) {
	PrintJewelHeader()
	for _, j := range jewels {
		PrintJewel(j)
	}
}
