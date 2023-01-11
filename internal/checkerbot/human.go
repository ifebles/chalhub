package checkerbot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

type human struct{}

var HumanPlayer = &human{}

func (h *human) HandlePlay(pl *player) Signal {
	var optPieces []*Piece

	if slayers := FilterSlayingOptions(*pl); len(slayers) > 0 {
		optPieces = slayers
	} else {
		optPieces = FilterSimpleOptions(*pl)
	}

	if len(optPieces) == 0 {
		canEnemyPlay := false

		if enemyMoves := FilterSlayingOptions(*pl.Enemy); len(enemyMoves) > 0 {
			canEnemyPlay = true
		} else if enemyMoves := FilterSimpleOptions(*pl.Enemy); len(enemyMoves) > 0 {
			canEnemyPlay = true
		}

		fmt.Println()
		modutil.PrintSystem("no moves can be made\n")

		if canEnemyPlay {
			util.PauseExecution()
			return Next
		} else {
			fmt.Print(" * Game ended in a truce *\n\n")
			return Quit
		}
	}

	playerPieceOption := getPlayerPieceOption(optPieces, pl.Char)
	selectedPiece, sgn := processSelectedPlayerPiece(playerPieceOption, optPieces)

	switch sgn {
	case Quit:
		return Quit // TODO: confirm quit

	case Print:
		return Print
	}

	trees := CreateTreeMaps(*pl, selectedPiece)
	plays := GetPiecePlays(trees)

	selectedPlayOption := handlePiecePlays(plays, PointToCoord(selectedPiece.Point), 10)

	// Move back:
	if selectedPlayOption == 0 {
		return Print
	}

	selectedPlay := plays[selectedPlayOption-1]

	modutil.PrintSystem("executing play...")
	ExecutePlay(pl, selectedPiece, selectedPlay)
	modutil.PrintSystem("done\n")

	return Cont
}

func readPlayerInput(attemptLimit int, pieces []*Piece) any {
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
			p := CoordToPoint(result)

			_, ok = util.Find(pieces, func(i *Piece) bool {
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

func getPlayerPieceOption(pieces []*Piece, plchr rune) any {
	fmt.Printf("\nSelect a piece: (player: %s)\n", string(plchr))
	fmt.Print("(P - print the board and options)\n\n")

	playerOptions := util.Map(pieces, func(i *Piece) string {
		suf := ""

		if i.IsKing {
			suf = " (king)"
		}

		return PointToCoord(i.Point) + suf
	})

	options := modutil.GetFormattedOptions(playerOptions, "Exit", 4)
	fmt.Printf("%s\n\n", strings.Join(options, "\n"))
	option := readPlayerInput(10, pieces)

	return option
}

func processSelectedPlayerPiece(opt any, pieces []*Piece) (*Piece, Signal) {
	var piece *Piece

	switch o := any(opt).(type) {
	case int:
		if o == 0 {
			return nil, Quit
		}

		piece = pieces[o-1]

	case string:
		if o == "P" {
			return nil, Print
		}

		piece, _ = util.Find(pieces, func(i *Piece) bool {
			p := CoordToPoint(o)
			return i.Point.X == p.X && i.Point.Y == p.Y
		})

	default:
		panic(fmt.Sprintf("unknown value (%T): %v", o, o))
	}

	return piece, Cont
}

func handlePiecePlays(plays []Play, pi string, attemptLimit int) int {
	fmt.Printf("\nSelect a play: (piece: %s)\n\n", pi)

	playerOptions := util.Map(plays, func(i Play) string {
		result := strings.Join(i.Breadcrumbs, ", ")
		slaycoords := util.Map(i.Slays, func(it *Point) string {
			return PointToCoord(*it)
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
