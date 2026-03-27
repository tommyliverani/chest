package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
)

func AskJewelAndRemove(jewelKind string) {
	chestJewels := factory.GetAllJewelsByKindFromOpenChests(jewelKind)
	if len(chestJewels) == 0 {
		fmt.Fprintf(os.Stderr, "No jewels of kind %s found in open chests\n", jewelKind)
		return
	}
	chestName, jewelToRemove := factory.SelectJewel("Select jewel to remove: ", chestJewels)
	removeJewelFromChest(chestName, jewelToRemove)
}

func RemoveJewelFromChestByName(jewelKind string, jewelName string) {

	chestJewels := factory.FindJewelsByKindAndNameFromOpenChests(jewelKind, jewelName)
	if len(chestJewels) == 0 {
		common.PrintErrorAndExit(fmt.Sprintf("No jewels of kind %s with name %s found in open chests", jewelKind, jewelName))
	}
	var chestName string
	var jewelToRemove factory.Jewel
	if len(chestJewels) > 1 {
		chestName, jewelToRemove = factory.SelectJewel("Select jewel to remove: ", chestJewels)
	} else {
		chestName = chestJewels[0].ChestName
		jewelToRemove = chestJewels[0].Jewel
	}
	removeJewelFromChest(chestName, jewelToRemove)
}

func removeJewelFromChest(chestName string, jewelToRemove factory.Jewel) {
	chest, found := factory.FindChestByName(chestName)
	if !found {
		common.PrintErrorAndExit(fmt.Sprintf("Chest %s not found", chestName))
	}
	keyJewel, found := factory.GetKeyJewelFromSession(chest.GetId())
	if !found {
		common.PrintErrorAndExit(fmt.Sprintf("Key jewel for chest %s not found in session", chestName))
	}
	confirm := common.SelectField(fmt.Sprintf("Are you sure you want to remove '%s'?", jewelToRemove.GetName()), []string{"No", "Yes"})
	if confirm != "Yes" {
		fmt.Println("Remove cancelled")
		return
	}
	err := chest.RemoveJewel(jewelToRemove, keyJewel)
	common.Check(err)
	_, err = factory.SaveOrUpdateChest(chest)
	common.Check(err)
	fmt.Printf("Jewel %s removed from chest %s\n", jewelToRemove.GetName(), chestName)
}
