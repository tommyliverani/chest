package command

import (
	"chest/internal/factory"
	"fmt"
)

func AskNameAndPrintJewel(jewelKind string) {
	chestJewels := factory.GetAllJewelsByKindFromOpenChests(jewelKind)
	if len(chestJewels) == 0 {
		fmt.Printf("No jewels of kind %s found in open chests\n", jewelKind)
		return
	}
	_, jewelToPrint := factory.SelectJewel("Select jewel to print: ", chestJewels)
	jewelToPrint.Print()
}

func PrintJewelByName(jewelKind string, jewelName string) {
	chestJewels := factory.FindJewelsByKindAndNameFromOpenChests(jewelKind, jewelName)
	if len(chestJewels) == 0 {
		fmt.Printf("No jewels of kind %s found in open chests\n", jewelKind)
		return
	}
	var jewelToPrint factory.Jewel
	if len(chestJewels) > 1 {
		_, jewelToPrint = factory.SelectJewel("Select jewel to print: ", chestJewels)
	} else {
		jewelToPrint = chestJewels[0].Jewel
	}
	jewelToPrint.Print()
}
