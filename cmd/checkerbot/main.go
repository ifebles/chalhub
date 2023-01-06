// Challenge inspiration from: https://www.hackerrank.com/challenges/checkers

package checkerbot

import (
	"fmt"
	"strings"

	"github.com/ifebles/chalhub/internal/checkerbot"
	"github.com/ifebles/chalhub/internal/modutil"
)

func GetName() string {
	return "Beat the checker-bot"
}

func Run() {
	fmt.Printf("** Welcome to the \"%s\" challenge **\n\n", GetName())

	defer checkerbot.Board.Clear()

	fmt.Print("Select the play mode:\n\n")

	options := modutil.GetFormattedOptions(checkerbot.PlayModes[:], "Exit", 3)
	fmt.Printf("%s\n\n", strings.Join(options, "\n"))

	selectedOption := modutil.GetIntOptionFromUser(10, len(checkerbot.PlayModes))

	if selectedOption == 0 {
		return
	}

	fmt.Println()
	modutil.PrintSystem("generating board...")

	// checkerbot.Board.Initialize()
	modutil.PrintSystem("done")

	fmt.Println()
	fmt.Println(checkerbot.Board)
}
