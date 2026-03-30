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

func ReadField(prompt string) string {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	Check(err)
	return strings.TrimSpace(input)
}

func SelectField(prompt string, fields []string) string {
	promptUI := promptui.Select{
		Label: prompt,
		Items: fields,
	}
	_, result, err := promptUI.Run()
	Check(err)
	return result
}

func SelectFieldWithIndex(prompt string, fields []string) (int, string) {
	promptUI := promptui.Select{
		Label: prompt,
		Items: fields,
	}
	index, result, err := promptUI.Run()
	Check(err)
	return index, result
}

func ReadSecret(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	Check(err)
	fmt.Println()
	return strings.TrimSpace(string(bytePassword))
}
