package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
)

func AskJewelAndEdit(kind string) {
	chestJewels := factory.GetAllJewelsByKindFromOpenChests(kind)
	if len(chestJewels) == 0 {
		fmt.Fprintf(os.Stderr, "No jewels of kind %s found in open chests\n", kind)
		return
	}
	chestName, jewelToEdit := factory.SelectJewel("Select jewel to edit: ", chestJewels)
	editJewelInChest(chestName, jewelToEdit, jewelToEdit.GetName())
}

func EditJewelByName(name string, kind string) {
	chestJewels := factory.FindJewelsByKindAndNameFromOpenChests(kind, name)
	if len(chestJewels) == 0 {
		common.PrintErrorAndExit(fmt.Sprintf("No jewels of kind %s with name %s found in open chests", kind, name))
	}
	var jewelToEdit factory.Jewel
	var chestName string
	if len(chestJewels) > 1 {
		chestName, jewelToEdit = factory.SelectJewel("Select jewel to edit: ", chestJewels)
	} else {
		chestName = chestJewels[0].ChestName
		jewelToEdit = chestJewels[0].Jewel
	}
	editJewelInChest(chestName, jewelToEdit, name)
}

func editJewelInChest(chestName string, jewelToEdit factory.Jewel, oldName string) {
	chest, found := factory.FindChestByName(chestName)
	if !found {
		common.PrintErrorAndExit(fmt.Sprintf("Chest %s not found", chestName))
	}
	keyJewel, found := factory.GetKeyJewelFromSession(chest.GetId())
	if !found {
		common.PrintErrorAndExit(fmt.Sprintf("Key jewel for chest %s not found in session", chestName))
	}
	common.Check(jewelToEdit.Edit())
	if jewelToEdit.GetName() != oldName {
		sameNameAndKindJewels := factory.FindJewelsByKindAndNameFromOpenChests(jewelToEdit.GetKind(), jewelToEdit.GetName())
		if len(sameNameAndKindJewels) > 0 {
			common.PrintErrorAndExit(fmt.Sprintf("A jewel of kind %s with name %s already exists in open chests", jewelToEdit.GetKind(), jewelToEdit.GetName()))
		}
	}
	err := chest.UpdateJewel(oldName, jewelToEdit, keyJewel)
	if err != nil {
		common.PrintErrorAndExit(fmt.Sprintf("Error updating jewel: %s", err.Error()))
	}
	_, err = factory.SaveOrUpdateChest(chest)
	common.Check(err)
	fmt.Printf("Jewel %s in chest %s edited\n", jewelToEdit.GetName(), chestName)
}
