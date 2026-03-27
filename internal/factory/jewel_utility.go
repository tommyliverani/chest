package factory

import (
	"chest/internal/common"
	"fmt"
	"sync"
)

type ChestJewel struct {
	ChestName string
	Jewel     Jewel
}

func waitAndCloseJewelChannel(wg *sync.WaitGroup, ch chan<- ChestJewel) {
	wg.Wait()
	close(ch)
}

func getJewelsFromChestId(id string, ch chan<- ChestJewel, wg *sync.WaitGroup) {
	defer wg.Done()
	chestJson, err := common.GetChestJsonById(id)
	if err != nil {
		return
	}
	chest, err := ParseChest(chestJson)
	if err != nil || chest == nil {
		return
	}
	keyJewel, found := GetKeyJewelFromSession(id)
	if !found {
		return
	}
	jewels, err := chest.GetJewels(keyJewel)
	if err != nil {
		return
	}
	for _, jewel := range jewels {
		ch <- ChestJewel{ChestName: chest.GetName(), Jewel: jewel}
	}
}

func GetAllJewelsFromOpenChests() []ChestJewel {
	ids := GetOpenChestIds()
	jewelsCh := make(chan ChestJewel)
	wg := sync.WaitGroup{}

	for _, id := range ids {
		wg.Add(1)
		go getJewelsFromChestId(id, jewelsCh, &wg)
	}
	go waitAndCloseJewelChannel(&wg, jewelsCh)
	var result []ChestJewel
	for cj := range jewelsCh {
		result = append(result, cj)
	}
	return result
}

func GetAllJewelsByKindFromOpenChests(kind string) []ChestJewel {
	all := GetAllJewelsFromOpenChests()
	var result []ChestJewel
	for _, cj := range all {
		if cj.Jewel.GetKind() == kind {
			result = append(result, cj)
		}
	}
	return result
}

func FindJewelsByKindAndNameFromOpenChests(kind string, name string) []ChestJewel {
	all := GetAllJewelsFromOpenChests()
	var result []ChestJewel
	for _, cj := range all {
		if cj.Jewel.GetKind() == kind && cj.Jewel.GetName() == name {
			result = append(result, cj)
		}
	}
	return result
}

func SelectJewel(prompt string, chestJewels []ChestJewel) (string, Jewel) {
	entryStrings := make([]string, len(chestJewels))
	for i, cj := range chestJewels {
		entryStrings[i] = GetJewelStringForSelect(cj.Jewel, cj.ChestName)
	}
	index, _ := common.SelectFieldWithIndex(prompt, entryStrings)
	if index >= 0 && index < len(chestJewels) {
		return chestJewels[index].ChestName, chestJewels[index].Jewel
	}
	return "", nil
}

func CheckJewelName(jewelName string, chestName string, chestJewels []ChestJewel) error {
	for _, cj := range chestJewels {
		if cj.Jewel.GetName() == jewelName && cj.ChestName == chestName {
			return nil
		}
	}
	return fmt.Errorf("Jewel with name %s not found in chest %s", jewelName, chestName)
}
