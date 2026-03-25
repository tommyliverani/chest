package command

import (
	"chest/internal/common"
	"fmt"
)

func RemoveJewelByName(kind string, name string) {

}

func RemoveJewel(kind string) {
	jewelName, err := common.ReadField("Insert jewel name: ")
	if err != nil {
		fmt.Printf("Error reading jewel name: %v\n", err)
		return
	}
	RemoveJewelByName(kind, jewelName)
}
