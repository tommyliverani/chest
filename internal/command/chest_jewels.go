package command

import (
	"chest/internal/factory"
	"fmt"
)

func ListJewels() {
	ids := factory.GetOpenChestIds()
	if len(ids) == 0 {
		fmt.Println("No open chests available")
		return
	}
	chestJewels := factory.GetAllJewelsFromOpenChests()
	if len(chestJewels) == 0 {
		fmt.Println("No jewels found in open chests")
		return
	}
	factory.PrintJewelHeader()
	for _, cj := range chestJewels {
		factory.PrintJewel(cj.Jewel, cj.ChestName)
	}
	factory.PrintJewelFooter(len(chestJewels))
}

func ListJewelsByKind(kind string) {
	ids := factory.GetOpenChestIds()
	if len(ids) == 0 {
		fmt.Println("No open chests available")
		return
	}
	chestJewels := factory.GetAllJewelsByKindFromOpenChests(kind)
	if len(chestJewels) == 0 {
		fmt.Println("No jewels found in open chests")
		return
	}
	factory.PrintJewelHeader()
	for _, cj := range chestJewels {
		factory.PrintJewel(cj.Jewel, cj.ChestName)
	}
	factory.PrintJewelFooter(len(chestJewels))
}
