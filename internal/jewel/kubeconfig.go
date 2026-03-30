package jewel

import (
	"bufio"
	"chest/internal/common"
	"chest/internal/factory"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const KUBECONFIG_KIND = "kube"

type Kubeconfig struct {
	baseJewel
	Kubeconfig string `json:"kubeconfig"`
}

func (k *Kubeconfig) GetEmoji() string { return "☸️ " }

func (k *Kubeconfig) ToJson() (json.RawMessage, error) { return json.Marshal(k) }

func ParseKubeconfig(data json.RawMessage) (*Kubeconfig, error) {
	var kc Kubeconfig
	if err := json.Unmarshal(data, &kc); err != nil {
		return nil, err
	}
	return &kc, nil
}

func CreateKubeconfig(name string, description string) (*Kubeconfig, error) {
	fmt.Println("Paste the kubeconfig content (press Ctrl+D when done):")
	config := readMultilineInput()
	return &Kubeconfig{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        KUBECONFIG_KIND,
			Description: description,
		},
		Kubeconfig: config,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(KUBECONFIG_KIND, func(name string, description string) (factory.Jewel, error) {
		return CreateKubeconfig(name, description)
	})
	factory.RegisterJewelParser(KUBECONFIG_KIND, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseKubeconfig(data)
	})
	factory.RegisterJewelHelp(KUBECONFIG_KIND, factory.JewelHelp{
		Emoji:    "☸️ ",
		Name:     "kube",
		Short:    "stores a kubeconfig and merges it into ~/.kube/config",
		Behavior: "merges the kubeconfig into ~/.kube/config and sets current-context",
		Operations: map[string]string{
			"add":   "saves a new kubeconfig (paste content, Ctrl+D to finish)",
			"ls":    "lists all kubeconfig entries in the open chest",
			"edit":  "edits name, description or kubeconfig content",
			"rm":    "removes a kubeconfig entry from the open chest",
			"print": "shows the raw kubeconfig YAML after confirmation",
			"copy":  "copies the kubeconfig YAML to the clipboard",
		},
	})
}

func (k *Kubeconfig) Edit() error {
	field := common.SelectField("which field do you want to edit?", []string{"name", "description", "config"})
	switch field {
	case "name":
		k.Name = common.ReadField("Insert new name: ")
	case "description":
		k.Description = common.ReadField("Insert new description: ")
	case "config":
		fmt.Println("Paste the new kubeconfig content (press Ctrl+D when done):")
		k.Kubeconfig = readMultilineInput()
	}
	return nil
}

func (k *Kubeconfig) Print() {
	fmt.Println(k.Kubeconfig)
}

func (k *Kubeconfig) Copy() {}

// Use merges the kubeconfig into ~/.kube/config and sets current-context.
func (k *Kubeconfig) Use() {
	kubeDir := filepath.Join(os.Getenv("HOME"), ".kube")
	if err := os.MkdirAll(kubeDir, 0700); err != nil { //nolint:gosec
		common.Check(err)
		return
	}
	kubeConfigPath := filepath.Join(kubeDir, "config")

	newCfg, err := parseKubeconfigYaml([]byte(k.Kubeconfig))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse kubeconfig: %v\n", err)
		os.Exit(1)
	}

	existing, err := loadExistingKubeconfig(kubeConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read existing kubeconfig: %v\n", err)
		os.Exit(1)
	}

	merged := mergeKubeconfigs(existing, newCfg)

	out, err := yaml.Marshal(merged)
	if err != nil {
		common.Check(err)
		return
	}
	if err := os.WriteFile(kubeConfigPath, out, 0600); err != nil { //nolint:gosec
		common.Check(err)
		return
	}
	fmt.Printf("Switched to context '%s'\n", merged["current-context"])
}

func readMultilineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}

func parseKubeconfigYaml(data []byte) (map[string]interface{}, error) {
	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadExistingKubeconfig(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path) //nolint:gosec
	if os.IsNotExist(err) {
		return map[string]interface{}{
			"apiVersion":      "v1",
			"kind":            "Config",
			"clusters":        []interface{}{},
			"users":           []interface{}{},
			"contexts":        []interface{}{},
			"current-context": "",
			"preferences":     map[string]interface{}{},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return parseKubeconfigYaml(data)
}

// mergeKubeconfigs merges src into dst: adds/replaces clusters, users, contexts by name,
// and sets current-context to src's current-context.
func mergeKubeconfigs(dst, src map[string]interface{}) map[string]interface{} {
	for _, key := range []string{"clusters", "users", "contexts"} {
		srcItems := toSlice(src[key])
		dstItems := toSlice(dst[key])
		for _, srcItem := range srcItems {
			srcMap := toMap(srcItem)
			name, _ := srcMap["name"].(string)
			replaced := false
			for i, dstItem := range dstItems {
				if toMap(dstItem)["name"] == name {
					dstItems[i] = srcItem
					replaced = true
					break
				}
			}
			if !replaced {
				dstItems = append(dstItems, srcItem)
			}
		}
		dst[key] = dstItems
	}
	if ctx, ok := src["current-context"]; ok {
		dst["current-context"] = ctx
	}
	return dst
}

func toSlice(v interface{}) []interface{} {
	if v == nil {
		return []interface{}{}
	}
	s, ok := v.([]interface{})
	if !ok {
		return []interface{}{}
	}
	return s
}

func toMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{}
}
