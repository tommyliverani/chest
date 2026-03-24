package command

import (
	. "chest/internal/common"
	. "chest/internal/jewel"
	"fmt"
)

func ReadJewelName(kind string) {
	name, err := ReadJewelName()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	description := ReadJewelDescription()
	jewel, err := CreateJewel(kind, name, description)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}

func AddJewel(kind string) {

}

func AddJewelByName(kind string, name string) {

}
