package main

import (
	_ "chest/internal/chest"
	"chest/internal/command"
	"chest/internal/factory"
	"fmt"
	"os"
	"slices"
)

// Version viene iniettata a tempo di compilazione tramite il flag:
// -ldflags "-X main.Version=$(VERSION)"
// Se non viene iniettata, il valore di default sarà "development"
var Version = "development"

var chestCommands = []string{"create", "ls", "rm", "edit", "jewels", "js", "open", "close"}
var jewelCommands = []string{"add", "rm", "ls", "edit", "jewels", "js", "copy", "print"}

func printHelp() {
	fmt.Printf("🏴‍☠️  Chest (Version: %s)\n\n", Version)
	fmt.Println("Available jewel kinds:")
	for _, kind := range factory.GetAvailableJewelKinds() {
		fmt.Println(factory.ShortHelp(kind))
	}
	fmt.Println()
	fmt.Printf("Chest commands: %v\n", chestCommands)
	fmt.Println("Tip: run './chest <kind> help' for details on a jewel type.")
}

func main() {

	if len(os.Args) < 2 || os.Args[1] == "help" {
		printHelp()
		return
	}

	if slices.Contains(chestCommands, os.Args[1]) {
		switchChestCommand(os.Args[1], os.Args)
		return
	}

	if slices.Contains(factory.GetAvailableJewelKinds(), os.Args[1]) {
		if len(os.Args) < 2 {
			fmt.Printf("Please provide a command for the jewel (available commands: %v) or a jewel kind(available jewel kinds: %v)\n", jewelCommands, factory.GetAvailableJewelKinds())
			return
		}

		if len(os.Args) == 3 && os.Args[2] == "help" {
			fmt.Print(factory.LongHelp(os.Args[1]))
			return
		}

		if len(os.Args) > 2 && slices.Contains(jewelCommands, os.Args[2]) {
			switchJewelCommand(os.Args[1], os.Args[2], os.Args)
			return
		}

		if len(os.Args) == 3 && slices.Contains(factory.GetAvailableJewelKinds(), os.Args[1]) {
			command.UseJewelByName(os.Args[1], os.Args[2])
			return
		}
		if len(os.Args) == 2 && slices.Contains(factory.GetAvailableJewelKinds(), os.Args[1]) {
			command.AskNameAndUseJewel(os.Args[1])
			return
		}

		fmt.Printf("Unknown jewel command: %s\n", os.Args[2])
		return
	}
	fmt.Printf("unknown command: %s\n", os.Args[1])

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
	case "open":
		if len(args) > 2 {
			command.OpenChestByName(args[2])
		} else {
			command.AskNameAndOpenChest()
		}
	case "close":
		if len(args) > 2 {
			command.CloseChestByName(args[2])
		} else {
			command.CloseChest()
		}
	case "jewels", "js":
		command.ListJewels()
	default:
		return
	}
}

func switchJewelCommand(kind string, jewelCommand string, args []string) {
	switch jewelCommand {
	case "add":
		if len(args) > 3 {
			command.AddJewelToChestByName(kind, args[3])
		} else {
			command.AskNameAndAddJewel(kind)
		}
	case "rm":
		if len(args) > 3 {
			command.RemoveJewelFromChestByName(kind, args[3])
		} else {
			command.AskJewelAndRemove(kind)
		}
	case "ls", "jewels", "js":
		command.ListJewelsByKind(kind)
	case "edit":
		if len(args) > 3 {
			command.EditJewelByName(args[3], kind)
		} else {
			command.AskJewelAndEdit(kind)
		}
	case "copy":
		if len(args) > 3 {
			command.CopyJewelByName(kind, args[3])
		} else {
			command.AskNameAndCopyJewel(kind)
		}
	case "print":
		if len(args) > 3 {
			command.PrintJewelByName(kind, args[3])
		} else {
			command.AskNameAndPrintJewel(kind)
		}

	default:
		fmt.Printf("Unknown jewel command: %s\n", jewelCommand)
	}
}
