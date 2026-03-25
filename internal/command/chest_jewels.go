package command

import (
	"chest/internal/factory"
	"fmt"
	"sync"
	"sync/atomic"
)

func ListJewelsWithKind(kind string) {
	openChests, err := factory.GetOpenChestNames()
	if err != nil {
		fmt.Printf("Error retrieving open chests: %v\n", err)
		return
	}
	if len(openChests) == 0 {
		fmt.Println("No chest opened. Please open a chest to list its jewels.")
		return
	}
	wg := sync.WaitGroup{}
	var jewelCount atomic.Int64
	factory.PrintJewelHeader()
	for _, chestName := range openChests {
		wg.Add(1)
		go ParseChestAndPrintJewelsWithKind(chestName, kind, &wg, &jewelCount)
	}
	wg.Wait()
	factory.PrintJewelFooter(int(jewelCount.Load()))
}

func ListJewels() {
	openChests, err := factory.GetOpenChestNames()
	if err != nil {
		fmt.Printf("Error retrieving open chests: %v\n", err)
		return
	}
	if len(openChests) == 0 {
		fmt.Println("No chest opened. Please open a chest to list its jewels.")
		return
	}
	wg := sync.WaitGroup{}
	var jewelCount atomic.Int64
	factory.PrintJewelHeader()
	for _, chestName := range openChests {
		wg.Add(1)
		go ParseChestAndPrintJewels(chestName, &wg, &jewelCount)
	}
	wg.Wait()
	factory.PrintJewelFooter(int(jewelCount.Load()))
}

func ParseChestAndPrintJewels(chestName string, wg *sync.WaitGroup, jewelCount *atomic.Int64) {
	defer wg.Done()
	defer func() {
		recover()
	}()
	chest, err := factory.GetExistingChest(chestName)
	keyJewel, found, err := factory.GetKeyJewel(chestName)
	if found == false {
		return
	}
	if err != nil {
		return
	}
	jewels, err := chest.GetJewels(keyJewel)
	if err != nil {
		return
	}
	for _, j := range jewels {
		factory.PrintJewel(j, chestName)
		jewelCount.Add(1)
	}
}

func ParseChestAndPrintJewelsWithKind(chestName string, kind string, wg *sync.WaitGroup, jewelCount *atomic.Int64) {
	defer wg.Done()
	defer func() {
		recover()
	}()
	chest, err := factory.GetExistingChest(chestName)
	keyJewel, found, err := factory.GetKeyJewel(chestName)
	if found == false {
		return
	}
	if err != nil {
		return
	}
	jewels, err := chest.GetJewels(keyJewel)
	if err != nil {
		return
	}
	for _, j := range jewels {
		if j.GetKind() == kind {
			factory.PrintJewel(j, chestName)
			jewelCount.Add(1)
		}
	}
}
