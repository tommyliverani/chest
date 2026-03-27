package factory

import (
	"chest/internal/common"
	"sync"
)

func waitAndCloseStringChannel(wg *sync.WaitGroup, ch chan<- string) {
	wg.Wait()
	close(ch)
}

func waitAndCloseChestChannel(wg *sync.WaitGroup, ch chan<- Chest) {
	wg.Wait()
	close(ch)
}

func parseChestName(id string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	chestJson, err := common.GetChestJsonById(id)
	if err != nil {
		return
	}
	name, err := common.GetNameFromJson(chestJson)
	if err != nil {
		return
	}
	ch <- name
}

func parseChestById(id string, ch chan<- Chest, wg *sync.WaitGroup) {
	defer wg.Done()
	chestJson, err := common.GetChestJsonById(id)
	if err != nil {
		return
	}
	chest, err := ParseChest(chestJson)
	if err != nil || chest == nil {
		return
	}
	ch <- chest
}

func GetExistingChestNames() []string {
	ids, err := common.GetExistingChestIds()
	common.Check(err)
	namesCh := make(chan string)
	wg := sync.WaitGroup{}
	for _, id := range ids {
		wg.Add(1)
		go parseChestName(id, namesCh, &wg)
	}
	go waitAndCloseStringChannel(&wg, namesCh)
	var names []string
	for name := range namesCh {
		names = append(names, name)
	}
	return names
}

func GetAllChests() []Chest {
	ids, err := common.GetExistingChestIds()
	common.Check(err)
	chestsCh := make(chan Chest)
	wg := sync.WaitGroup{}
	for _, id := range ids {
		wg.Add(1)
		go parseChestById(id, chestsCh, &wg)
	}
	go waitAndCloseChestChannel(&wg, chestsCh)
	var chests []Chest
	for chest := range chestsCh {
		chests = append(chests, chest)
	}
	return chests
}

func SelectChest(prompt string, chests []Chest) Chest {
	entrys := make([]string, 0, len(chests))
	for _, chest := range chests {
		entrys = append(entrys, GetChestStringForSelect(chest))
	}
	index, _ := common.SelectFieldWithIndex(prompt, entrys)
	return chests[index]
}

func FindChestByName(name string) (Chest, bool) {
	chests := GetAllChests()
	for _, chest := range chests {
		if chest.GetName() == name {
			return chest, true
		}
	}
	return nil, false
}

func GetAllOpenChests() []Chest {
	ids := GetOpenChestIds()
	chestsCh := make(chan Chest)
	wg := sync.WaitGroup{}
	for _, id := range ids {
		wg.Add(1)
		go parseChestById(id, chestsCh, &wg)
	}
	go waitAndCloseChestChannel(&wg, chestsCh)
	var chests []Chest
	for chest := range chestsCh {
		chests = append(chests, chest)
	}
	return chests
}
