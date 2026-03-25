package command

import (
	"chest/internal/common"
	"chest/internal/factory"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

func ListChests() {

	chestNames, err := common.GetExistingChestNames()
	if len(chestNames) == 0 {
		fmt.Println("No chest available")
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving chest: %v\n", err)
		os.Exit(1)
	}
	factory.PrintChestHeader()
	wg := sync.WaitGroup{}
	var chestCount atomic.Int64
	for _, chestName := range chestNames {
		wg.Add(1)
		go ParseAndPrintChest(chestName, &wg, &chestCount)
	}
	wg.Wait()
	factory.PrintChestFooter(int(chestCount.Load()))
}

func ParseAndPrintChest(chestName string, wg *sync.WaitGroup, chestCount *atomic.Int64) {
	defer wg.Done()
	defer func() {
		recover()
	}()
	chest, err := factory.GetExistingChest(chestName)
	if err == nil {
		factory.PrintChest(chest)
		chestCount.Add(1)
	}
}
