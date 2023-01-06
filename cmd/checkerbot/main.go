// Challenge inspiration from: https://www.hackerrank.com/challenges/checkers

package checkerbot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ifebles/chalhub/internal/checkerbot"
	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

type signal int

const (
	cont signal = iota
	print
	quit
)

func GetName() string {
	return "Beat the checker-bot"
}

func Run() {
	fmt.Printf("** Welcome to the \"%s\" challenge **\n\n", GetName())

	defer checkerbot.Board.Clear()

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

	currentPlayer := checkerbot.StartGame(checkerbot.PlayMode(selectedOption - 1))

	for {
		fmt.Printf("\nTurn #%d\n", checkerbot.GetTurnNumber())
		playerFlag, isAi := "", currentPlayer.Type == checkerbot.Ai

		if isAi {
			playerFlag = fmt.Sprintf("%s (%s)", string(currentPlayer.Char), strings.ToUpper(string(currentPlayer.Type)))
		} else {
			playerFlag = string(currentPlayer.Char)
		}

		fmt.Printf("Current player: %s\n\n", playerFlag)

		fmt.Println()
		fmt.Println(checkerbot.Board)

		if !isAi {
			playerPieces := checkerbot.GetPieces(currentPlayer)
			playerPieceOption := getPlayerPieceOption(playerPieces)
			selectedPiece, sgn := processSelectedPlayerPiece(playerPieceOption, playerPieces)

			switch sgn {
			case quit:
				return // TODO: confirm quit

			case print:
				continue
			}

			fmt.Printf("Selected piece: %v\n", selectedPiece)
		}

		util.PauseExecution()

		currentPlayer = checkerbot.EndTurn()
	}

}

func readPlayerInput(attemptLimit, valueLimit int) any {
	if attemptLimit <= 0 {
		panic("invalid limit given")
	}

	for x := 0; x < attemptLimit; x++ {
		result, err := util.ReadString(">> ")

		if err != nil {
			fmt.Println("invalid input")
			continue
		}

		num, err := strconv.Atoi(result)

		if err == nil {
			if min, max := 0, valueLimit; num < min || num > max {
				modutil.PrintAdvice("if an integer is provided, must be between %d and %d", min, max)
				continue
			}

			return num
		}

		result = strings.ToUpper(result)

		if result == "P" {
			return result
		}

		ok := false

		func() {
			defer func() {
				if err := recover(); err == nil {
					ok = true
				}
			}()

			checkerbot.CoordToPoint(result)
		}()

		if ok {
			return result
		}
	}

	return "P"
}

func getPlayerPieceOption(pieces []checkerbot.Piece) any {
	fmt.Print("\nSelect a piece:\n")
	fmt.Print("(P - print the board and options)\n\n")

	playerOptions := util.Map(pieces, func(i checkerbot.Piece) string {
		suf := ""

		if i.IsKing {
			suf = " (king)"
		}

		return checkerbot.PointToCoord(i.Point) + suf
	})

	options := modutil.GetFormattedOptions(playerOptions, "Exit", 4)

	fmt.Printf("%s\n\n", strings.Join(options, "\n"))

	option := readPlayerInput(10, len(pieces))

	return option
}

func processSelectedPlayerPiece(opt any, pieces []checkerbot.Piece) (*checkerbot.Piece, signal) {
	var piece *checkerbot.Piece

	switch o := any(opt).(type) {
	case int:
		if o == 0 {
			return nil, quit
		}

		piece = &pieces[o-1]

	case string:
		if o == "P" {
			return nil, print
		}

		piece, _ = util.Find(pieces, func(i checkerbot.Piece) bool {
			p := checkerbot.CoordToPoint(o)
			return i.Point.X == p.X && i.Point.Y == p.Y
		})

	default:
		panic(fmt.Sprintf("unknown value (%T): %v", o, o))
	}

	return piece, cont
}
