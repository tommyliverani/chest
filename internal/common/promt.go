package common

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/manifoldco/promptui"
)

func ReadChestName() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter chest name: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	trimmedInput := strings.TrimSpace(input)
	if slices.Contains(GetChestNames(), trimmedInput) {
		return "", fmt.Errorf("chest with name '%s' already exists", trimmedInput)
	}
	return trimmedInput, nil
}

func ReadChestDescription() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter chest description: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	trimmedInput := strings.TrimSpace(text)
	return trimmedInput, nil
}

func SelectChestKind(kinds []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select chest kind",
		Items: kinds,
	}
	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return result, nil
}

func SelectChest(chests []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select chest",
		Items: chests,
	}
	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return result, nil
}

func SelectChestField(fields []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select chest field",
		Items: fields,
	}
	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return result, nil

}
