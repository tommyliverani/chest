package main

import (
	. "chest/internal/chest"
	. "chest/internal/command"
	. "chest/internal/common"
	"fmt"
	"os"
)

// Version viene iniettata a tempo di compilazione tramite il flag:
// -ldflags "-X main.Version=$(VERSION)"
// Se non viene iniettata, il valore di default sarà "development"
var Version = "development"

func main() {
	EnsureDir(ChestBasePath)

	if len(os.Args) > 1 && os.Args[1] == "create" {
		CreateChest()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "ls" {
		ListChests()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "delete" {
		DeleteChest()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "edit" {
		EditChest()
		return
	}

	fmt.Printf("Chest (Version: %s) is running...\n", Version)
	fmt.Printf("Supported Chest \n")
	for _, name := range GetKinds() {
		fmt.Printf("%s\n", name)
	}

}

func EnsureDir(dir string) error {

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	return nil
}
