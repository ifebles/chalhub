package main

import (
	"fmt"
	"strings"

	"github.com/ifebles/chalhub/cmd/botsaver"
	"github.com/ifebles/chalhub/cmd/checkerbot"
	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

type challenge struct {
	name string
	run  func()
}

var challenges = []challenge{
	{botsaver.GetName(), botsaver.Run},
	{checkerbot.GetName(), checkerbot.Run},
}

func main() {
	fmt.Print("** Welcome to the Challenge HUB **\n")
	fmt.Print("* Here you can select which challenge to open *\n\n")

	for {
		fmt.Print("Select an option:\n\n")

		challengeNames := util.MapCollection(challenges, func(item challenge) string {
			return item.name
		})

		options := modutil.GetFormattedOptions(challengeNames)
		fmt.Printf("%s\n\n", strings.Join(options, "\n"))

		selectedOption := modutil.GetIntOptionFromUser(10, len(challenges))

		if selectedOption == 0 {
			fmt.Print("\nGoodbye!\n")
			return
		}

		////

		fmt.Println()
		modutil.PrintSystem("starting challenge...\n")

		challenges[selectedOption-1].run()

		fmt.Println()
		util.PauseExecution()
		fmt.Println()
		modutil.PrintSystem("closing challenge \"%s\"...", challenges[selectedOption-1].name)
		modutil.PrintSystem("returning to main HUB...\n")
	}
}
