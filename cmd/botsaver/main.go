// Challenge inspiration from: https://www.hackerrank.com/challenges/saveprincess

package botsaver

import (
	"fmt"

	"github.com/ifebles/chalhub/internal/botsaver"
	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

func GetName() string {
	return "Save the Princess"
}

func Run() {
	fmt.Printf("** Welcome to the \"%s\" solver **\n\n", GetName())

	defer botsaver.Grid.Clear()
	defer botsaver.Bot.Clear()

	gridSize := botsaver.Grid.GetGridSize(10)

	if gridSize == -1 {
		panic("unable to get a valid input")
	}

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
