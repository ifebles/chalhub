package checkerbot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

type Signal int

const (
	Cont Signal = iota
	Next
	Print
	Quit
)

type ai struct{}

type piecePlay struct {
	piece *Piece
	play  Play
}

var AI = &ai{}

func (ai *ai) DecidePlay(pl *player) Signal {
	var optPieces []*Piece
	slays := false

	fmt.Println()
	modutil.PrintSystem("analyzing plays...")

	if slayers := FilterSlayingOptions(*pl); len(slayers) > 0 {
		slays = true
		optPieces = slayers
	} else {
		optPieces = FilterSimpleOptions(*pl)
	}

	////

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

	////

	modutil.PrintSystem("selecting best option...")

	piece, play := selectPlay(pl, optPieces, slays)
	modutil.PrintSystem("executing play...")

	pieceCoord := PointToCoord(piece.Point)
	fmt.Printf("\nPlaying piece %s to %s\n\n", pieceCoord, PointToCoord(play.Dest))
	fmt.Print("Details:\n\n")
	fmt.Printf("- Piece: %s\n", pieceCoord)
	fmt.Printf("- Moves: %s\n", strings.Join(play.Breadcrumbs, ", "))

	if slays {
		coords := util.Map(play.Slays, func(i *Point) string { return PointToCoord(*i) })
		fmt.Printf("- Slays: %s\n", strings.Join(coords, ", "))
	}

	if play.IsKing && !piece.IsKing {
		fmt.Print("- Status: crowned\n")
	}

	////

	ExecutePlay(pl, piece, play)

	fmt.Println()
	modutil.PrintSystem("done\n")

	return Cont
}

func selectPlay(pl *player, pieces []*Piece, isSlay bool) (*Piece, Play) {
	rand.Seed(time.Now().UnixNano())

	// TODO: establish a point system for movements
	// TODO: if at least one crowned, pursuit enemy pieces avoiding risks

	if !isSlay {
		// TODO: check first if the next enemy move causes a slay
		// TODO: check if next move causes crowning
		selectedPiece := pieces[rand.Intn(len(pieces))]
		trees := CreateTreeMaps(*pl, selectedPiece)
		plays := GetPiecePlays(trees)

		return selectedPiece, plays[rand.Intn(len(plays))]
	}

	////

	potentialPlays := make(map[int][]piecePlay)
	maxSlay := 0

	for _, a := range pieces {
		trees := CreateTreeMaps(*pl, a)
		plays := GetPiecePlays(trees)

		slayCounts := util.Map(plays, func(i Play) int { return len(i.Slays) })
		localMaxSlay := slayCounts[0]
		slayMap := make(map[int][]int)

		////

		for x, a := range slayCounts {
			if a > localMaxSlay {
				localMaxSlay = a
			}

			if _, ok := slayMap[a]; ok {
				slayMap[a] = append(slayMap[a], x)
			} else {
				slayMap[a] = []int{x}
			}
		}

		if maxSlay < localMaxSlay {
			maxSlay = localMaxSlay
		}

		if localMaxSlay < maxSlay {
			continue
		}

		////

		// TODO: check first if the next enemy move causes a slay
		// TODO: check if next move causes crowning
		selectedIndex := slayMap[localMaxSlay][rand.Intn(len(slayMap[localMaxSlay]))]

		if _, ok := potentialPlays[localMaxSlay]; ok {
			potentialPlays[localMaxSlay] = append(potentialPlays[localMaxSlay], piecePlay{a, plays[selectedIndex]})
		} else {
			potentialPlays[localMaxSlay] = []piecePlay{{a, plays[selectedIndex]}}
		}
	}

	////

	if len(potentialPlays[maxSlay]) == 0 {
		return potentialPlays[maxSlay][0].piece, potentialPlays[maxSlay][0].play
	}

	// TODO: check first if the next enemy move causes a slay
	// TODO: check if next move causes crowning
	selectedInx := rand.Intn(len(potentialPlays[maxSlay]))
	return potentialPlays[maxSlay][selectedInx].piece, potentialPlays[maxSlay][selectedInx].play
}
