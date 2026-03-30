package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
)

func AskNameAndAddJewel(jewelKind string) {
	openChests := factory.GetOpenChestIds()
	if len(openChests) == 0 {
		fmt.Fprintln(os.Stderr, "No open chests available, open a chest before adding jewels")
		return
	}
	var chestToAddJewel factory.Chest
	if len(openChests) > 1 {
		fmt.Println("Multiple open chests available")
		openChests := factory.GetAllOpenChests()
		chestToAddJewel = factory.SelectChest("Select chest to add jewel: ", openChests)
	} else {
		chestToAddJewel, _ = factory.GetChestById(openChests[0])
	}
	name := common.ReadField("Insert jewel name: ")
	addJewelToChest(chestToAddJewel, jewelKind, name)
}

func AddJewelToChestByName(jewelKind string, jewelName string) {
	openChests := factory.GetOpenChestIds()
	if len(openChests) == 0 {
		fmt.Fprintln(os.Stderr, "No open chests available, open a chest before adding jewels")
		return
	}
	var chestToAddJewel factory.Chest
	if len(openChests) > 1 {
		fmt.Println("Multiple open chests available")
		openChests := factory.GetAllOpenChests()
		chestToAddJewel = factory.SelectChest("Select chest to add jewel: ", openChests)
	} else {
		chestToAddJewel, _ = factory.GetChestById(openChests[0])
	}
	addJewelToChest(chestToAddJewel, jewelKind, jewelName)
}

func addJewelToChest(chestToAddJewel factory.Chest, jewelKind string, name string) {
	sameNameAndKindJewels := factory.FindJewelsByKindAndNameFromOpenChests(jewelKind, name)
	if len(sameNameAndKindJewels) > 0 {
		common.PrintErrorAndExit(fmt.Sprintf("A jewel of kind %s with name %s already exists in open chests", jewelKind, name))
	}

	description := common.ReadField("Insert jewel description: ")
	jewel, err := factory.CreateJewel(jewelKind, name, description)
	common.Check(err)
	keyJewel, found := factory.GetKeyJewelFromSession(chestToAddJewel.GetId())
	if !found {
		common.PrintErrorAndExit("Key jewel not found in session")
	}

	err = chestToAddJewel.AddJewel(jewel, keyJewel)
	common.Check(err)
	_, err = factory.SaveOrUpdateChest(chestToAddJewel)
	common.Check(err)
	fmt.Printf("Jewel %s added to chest %s\n", name, chestToAddJewel.GetName())
}
