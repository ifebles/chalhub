// Challenge inspiration from: https://www.hackerrank.com/challenges/saveprincess

package main

import (
	"fmt"

	"chalhub/internal/botsaver"
	"chalhub/internal/modutil"
	"chalhub/pkg/util"
)

func main() {
	defer botsaver.Grid.Clear()
	fmt.Print("** Welcome to the \"Save the Princess\" solver **\n\n")

	gridSize := botsaver.Grid.GetGridSize()

	modutil.PrintSystem("Generating %dx%d grid...", gridSize, gridSize)
	ok := botsaver.Grid.GenerateGrid()

	if ok {
		modutil.PrintSystem("Done\n")
	} else {
		panic("Unable to generate grid")
	}

	util.PauseExecution()

	fmt.Println(botsaver.Grid)
}
