package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

var reader = bufio.NewReader(os.Stdin)

func ReadField(prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	trimmedInput := strings.TrimSpace(input)
	return trimmedInput, nil
}

func SelectField(prompt string, fields []string) (string, error) {
	promptUI := promptui.Select{
		Label: prompt,
		Items: fields,
	}
	_, result, err := promptUI.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func SelectFieldWithIndex(prompt string, fields []string) (int, string, error) {
	promptUI := promptui.Select{
		Label: prompt,
		Items: fields,
	}
	index, result, err := promptUI.Run()
	if err != nil {
		return -1, "", err
	}
	return index, result, nil
}

func ReadSecret(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return strings.TrimSpace(string(bytePassword)), nil
}
