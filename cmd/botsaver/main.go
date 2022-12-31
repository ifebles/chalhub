// Challenge inspiration from: https://www.hackerrank.com/challenges/saveprincess

package main

import (
	"fmt"

	"github.com/ifebles/chalhub/internal/botsaver"
	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

func main() {
	fmt.Print("** Welcome to the \"Save the Princess\" solver **\n\n")

	defer botsaver.Grid.Clear()
	defer botsaver.Bot.Clear()

	gridSize := botsaver.Grid.GetGridSize()

	modutil.PrintSystem("generating %dx%d grid...", gridSize, gridSize)
	ok := botsaver.Grid.GenerateGrid()

	if ok {
		modutil.PrintSystem("done\n")
	} else {
		panic("unable to generate grid")
	}

	util.PauseExecution()

	fmt.Printf("\n%v\n\n", botsaver.Grid)

	modutil.PrintSystem("calculating best route for the bot (m)...")

	botsaver.Bot.Feed(fmt.Sprint(botsaver.Grid))
	botsaver.Bot.DisplayPathToPrincess(func() {
		modutil.PrintSystem("done\n")
		util.PauseExecution()
	})
}
