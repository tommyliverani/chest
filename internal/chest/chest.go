package chest

import (
	. "chest/internal/common"
	"encoding/json"
	"fmt"
	"os"
)

type Chest interface {
	GetName() string
	GetKind() string
	GetDescription() string
	GetEmoji() string
	Delete() error
	Edit() error
}

func SaveChestToFile(filename string, chest Chest) error {
	data, err := json.MarshalIndent(chest, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func PrintChestHeader() {
	fmt.Println("\tNAME\t\tKIND\t\tDESCRIPTION")
}
func PrintChest(chest Chest) {
	fmt.Printf("%s\t%s\t\t%s\t\t%s\n", chest.GetEmoji(), chest.GetName(), chest.GetKind(), chest.GetDescription())
}

// da rivedere con go routines
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

func GetChestByName(name string) Chest {
	chestPath := ChestBasePath + "/" + name + ".json"
	chestFile, errRead := os.ReadFile(chestPath)
	if errRead == nil {
		chest, errParse := ParseChest(chestFile)
		if errParse == nil {
			return chest
		} else {
			fmt.Printf("failed to parse chest file %s: %v\n", chestPath, errParse)
		}
	} else {
		fmt.Printf("failed to read chest file %s: %v\n", chestPath, errRead)
	}
	return nil
}
