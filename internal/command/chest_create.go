package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
)

func CreateChest() {
	name := common.ReadField("Insert chest name: ")
	common.Check(factory.CheckChestName("", name))
	CreateChestByName(name)
}

func CreateChestByName(name string) {
	common.Check(factory.CheckChestName("", name))
	// kind := common.SelectField("Select chest kind: ", factory.GetAvailableChestKinds())
	kind := "aes" // aes is the only available chest kind for now, so we skip the selection
	description := common.ReadField("Insert chest description: ")
	chest, err := factory.CreateChest(kind, name, description)
	common.Check(err)
	chestPath, err := factory.SaveOrUpdateChest(chest)
	common.Check(err)
	fmt.Printf("Chest created in %s\n", chestPath)
}
