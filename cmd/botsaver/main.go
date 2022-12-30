// Challenge inspiration from: https://www.hackerrank.com/challenges/saveprincess

package main

import (
	"fmt"

	"chalhub/internal/botsaver"
	"chalhub/internal/modutil"
	"chalhub/pkg/util"
)

func main() {
	fmt.Print("** Welcome to the \"Save the Princess\" solver **\n\n")

	gridSize := botsaver.GetGridSize()

	modutil.PrintSystem("Generating %dx%d grid...", gridSize, gridSize)
	grid := botsaver.GenerateGrid(gridSize)
	modutil.PrintSystem("Done\n")

	util.PauseExecution()

	fmt.Println(grid)
}
