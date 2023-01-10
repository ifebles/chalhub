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
			var optPieces []*checkerbot.Piece

			if slayers := checkerbot.FilterSlayingOptions(*currentPlayer); len(slayers) > 0 {
				optPieces = slayers
			} else {
				optPieces = checkerbot.FilterSimpleOptions(*currentPlayer)
			}

			if len(optPieces) == 0 {
				canEnemyPlay := false

				if enemyMoves := checkerbot.FilterSlayingOptions(*currentPlayer.Enemy); len(enemyMoves) > 0 {
					canEnemyPlay = true
				} else if enemyMoves := checkerbot.FilterSimpleOptions(*currentPlayer.Enemy); len(enemyMoves) > 0 {
					canEnemyPlay = true
				}

				fmt.Println()
				modutil.PrintSystem("no moves can be made\n")

				if canEnemyPlay {
					util.PauseExecution()
					currentPlayer = checkerbot.EndTurn()

					continue
				} else {
					fmt.Print(" * Game ended in a truce *\n\n")
					return
				}
			}

			playerPieceOption := getPlayerPieceOption(optPieces, currentPlayer.Char)
			selectedPiece, sgn := processSelectedPlayerPiece(playerPieceOption, optPieces)

			switch sgn {
			case quit:
				return // TODO: confirm quit

			case print:
				continue
			}

			trees := checkerbot.CreateTreeMaps(*currentPlayer, selectedPiece)
			plays := checkerbot.GetPiecePlays(trees)

			selectedPlayOption := handlePiecePlays(plays, checkerbot.PointToCoord(selectedPiece.Point), 10)

			// Move back:
			if selectedPlayOption == 0 {
				continue
			}

			selectedPlay := plays[selectedPlayOption-1]

			modutil.PrintSystem("executing play...")
			checkerbot.ExecutePlay(currentPlayer, selectedPiece, selectedPlay)
			modutil.PrintSystem("done\n")
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

func readPlayerInput(attemptLimit int, pieces []*checkerbot.Piece) any {
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
			if min, max := 0, len(pieces); num < min || num > max {
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
			defer func() { recover() }()
			p := checkerbot.CoordToPoint(result)

			_, ok = util.Find(pieces, func(i *checkerbot.Piece) bool {
				return i.Point.X == p.X && i.Point.Y == p.Y
			})
		}()

		if ok {
			return result
		}
		// else
		modutil.PrintAdvice("input a valid coordinate/command")
	}

	return "P"
}

func getPlayerPieceOption(pieces []*checkerbot.Piece, plchr rune) any {
	fmt.Printf("\nSelect a piece: (player: %s)\n", string(plchr))
	fmt.Print("(P - print the board and options)\n\n")

	playerOptions := util.Map(pieces, func(i *checkerbot.Piece) string {
		suf := ""

		if i.IsKing {
			suf = " (king)"
		}

		return checkerbot.PointToCoord(i.Point) + suf
	})

	options := modutil.GetFormattedOptions(playerOptions, "Exit", 4)
	fmt.Printf("%s\n\n", strings.Join(options, "\n"))
	option := readPlayerInput(10, pieces)

	return option
}

func processSelectedPlayerPiece(opt any, pieces []*checkerbot.Piece) (*checkerbot.Piece, signal) {
	var piece *checkerbot.Piece

	switch o := any(opt).(type) {
	case int:
		if o == 0 {
			return nil, quit
		}

		piece = pieces[o-1]

	case string:
		if o == "P" {
			return nil, print
		}

		piece, _ = util.Find(pieces, func(i *checkerbot.Piece) bool {
			p := checkerbot.CoordToPoint(o)
			return i.Point.X == p.X && i.Point.Y == p.Y
		})

	default:
		panic(fmt.Sprintf("unknown value (%T): %v", o, o))
	}

	return piece, cont
}

func handlePiecePlays(plays []checkerbot.Play, pi string, attemptLimit int) int {
	fmt.Printf("\nSelect a play: (piece: %s)\n\n", pi)

	playerOptions := util.Map(plays, func(i checkerbot.Play) string {
		result := strings.Join(i.Breadcrumbs, ", ")
		slaycoords := util.Map(i.Slays, func(it *checkerbot.Point) string {
			return checkerbot.PointToCoord(*it)
		})

		if len(slaycoords) > 0 {
			result += fmt.Sprintf(" (Slays: %s)", strings.Join(slaycoords, ", "))
		}

		return result
	})

	options := modutil.GetFormattedOptions(playerOptions, "Back", 1)
	fmt.Printf("%s\n\n", strings.Join(options, "\n"))

	for x := 0; x < attemptLimit; x++ {
		resp, err := util.ReadInteger(">> ")

		if err != nil {
			modutil.PrintAdvice("an integer was expected")
			continue
		}

		if min, max := 0, len(plays); resp < min || resp > max {
			modutil.PrintAdvice("the number must be between %d and %d", min, max)
			continue
		}

		return resp
	}

	return 0
}
