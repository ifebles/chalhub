// Challenge inspiration from: https://www.hackerrank.com/challenges/checkers

package checkerbot

import (
	"fmt"
	"strings"

	"github.com/ifebles/chalhub/internal/checkerbot"
	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

func GetName() string {
	return "Beat the checker-bot"
}

func Run() {
	fmt.Printf("** Welcome to the \"%s\" challenge **\n\n", GetName())

	board := checkerbot.GetNewBoard()
	defer board.Clear()

	fmt.Print("Select the play mode:\n\n")

	playModes := util.Map(checkerbot.PlayModes[:], func(i checkerbot.PlayMode) string {
		return fmt.Sprint(i)
	})

	options := modutil.GetFormattedOptions(playModes, "Exit", 3)
	fmt.Printf("%s\n\n", strings.Join(options, "\n"))

	selectedOption := modutil.GetIntOptionFromUser(10, len(playModes))

	if selectedOption == 0 {
		return
	}

	fmt.Println()
	modutil.PrintSystem("starting game...")

	currentPlayer := checkerbot.StartGame(board, checkerbot.PlayMode(selectedOption-1))

	for {
		turn := checkerbot.GetTurnNumber()
		fmt.Printf("\nTurn #%d\n", turn)
		playerFlag, isAi := "", currentPlayer.Type == checkerbot.Ai

		if isAi {
			playerFlag = fmt.Sprintf("%s (%s)", string(currentPlayer.Char), strings.ToUpper(string(currentPlayer.Type)))
		} else {
			playerFlag = string(currentPlayer.Char)
		}

		fmt.Printf("Current player: %s\n\n", playerFlag)

		fmt.Println()
		fmt.Println(board)

		var sgn checkerbot.Signal

		if isAi {
			if turn == 1 {
				fmt.Println()
				util.PauseExecution()
			}

			sgn = checkerbot.AI.DecidePlay(currentPlayer)
		} else {
			sgn = checkerbot.HumanPlayer.HandlePlay(currentPlayer)
		}

		switch sgn {
		case checkerbot.Next:
			currentPlayer = checkerbot.EndTurn()
			continue

		case checkerbot.Print:
			continue

		case checkerbot.Quit:
			return
		}

		if len(*currentPlayer.Enemy.Pieces) == 0 {
			fmt.Print("\n ** The game has ended! **\n\n")
			fmt.Printf(" * Winner: %s *\n\n", playerFlag)

			return
		}

		util.PauseExecution()

		currentPlayer = checkerbot.EndTurn()
	}
}
