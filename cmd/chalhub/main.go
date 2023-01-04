package main

import (
	"fmt"
	"strings"

	"github.com/ifebles/chalhub/cmd/botsaver"
	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

type challenge struct {
	name string
	run  func()
}

var challenges = []challenge{
	{botsaver.GetName(), botsaver.Run},
}

func main() {
	fmt.Print("** Welcome to the Challenge HUB **\n")
	fmt.Print("* Here you can select which challenge to open *\n\n")

	for {
		fmt.Print("Select an option:\n\n")
		options := getOptions()

		fmt.Printf("%s\n\n", strings.Join(options, "\n"))

		selectedOption := getOptionFromUser(10)

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

func getOptions() []string {
	result := make([]string, 0, len(challenges)+1)

	for x := range challenges {
		result = append(result, fmt.Sprintf("\t%d) %s", x+1, challenges[x].name))
	}

	result = append(result, "\t0) Exit")

	return result
}

func getOptionFromUser(attemptLimit int) int {
	if attemptLimit <= 0 {
		panic("invalid limit given")
	}

	for x := 0; x < attemptLimit; x++ {
		result, err := util.ReadInteger(">> ")

		if err != nil {
			modutil.PrintAdvice("an integer was expected")
			continue
		}

		if min, max := 0, len(challenges); result < min || result > max {
			modutil.PrintAdvice("an integer between %d and %d was expected", min, max)
			continue
		}

		return result
	}

	return 0
}
