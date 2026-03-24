package main

import (
	"chest/internal/command"
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
	"slices"
)

// Version viene iniettata a tempo di compilazione tramite il flag:
// -ldflags "-X main.Version=$(VERSION)"
// Se non viene iniettata, il valore di default sarà "development"
var Version = "development"

func main() {
	if err := ensureDir(common.GetChestHome()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Printf("Chest (Version: %s) is running...\n", Version)
		return
	}

	if slices.Contains(command.ChestCommands, os.Args[1]) {
		switchChestCommand(os.Args[1], os.Args)
		return
	}

	if slices.Contains(factory.GetAvailableJewelKinds(), os.Args[1]) {
		if slices.Contains(command.JewelCommands, os.Args[2]) || len(os.Args) < 3 {
			switchJewelCommand(os.Args[2], os.Args)
			return
		}
		command.UseJewel(os.Args[1], os.Args[2])
		return
	}

}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func switchChestCommand(chestCommand string, args []string) {
	switch chestCommand {
	case "create":
		if len(args) > 2 {
			command.CreateChestByName(args[2])
		} else {
			command.CreateChest()
		}
	case "ls":
		command.ListChests()
	case "rm":
		if len(args) > 2 {
			command.DeleteChestByName(args[2])
		} else {
			command.DeleteChest()
		}
	case "edit":
		if len(args) > 2 {
			command.EditChestByName(args[2])
		} else {
			command.EditChest()
		}
	default:
		fmt.Printf("Unknown command: %s\n", chestCommand)
	}
}

func switchJewelCommand(jewelCommand string, args []string) {
	switch jewelCommand {
	case "add":
		if len(args) > 2 {
			command.AddJewelByName(args[1], args[4])
		} else {
			command.AddJewel(args[1])
		}
	case "ls":
		command.ListJewels(args[1])
	case "rm":
		if len(args) > 2 {
			command.RemoveJewelByName(args[2], args[4])
		} else {
			command.RemoveJewel(args[2])
		}
	case "edit":
		if len(args) > 2 {
			command.EditJewelByName(args[2], args[4])
		} else {
			command.EditJewel(args[2])
		}
	default:
		fmt.Printf("Unknown command: %s\n", jewelCommand)
	}
}
