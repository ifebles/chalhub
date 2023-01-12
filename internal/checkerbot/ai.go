package checkerbot

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

type playScore int

const (
	chasingMove playScore = 1 << (iota + 1)
	simpleSlayMove
	crowningMove
	kingSlayMove
	nullifyingMove
)

type ai struct{}

type piecePlay struct {
	piece *Piece
	play  Play
	score float64
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

	bestPlays := []piecePlay{}

	for _, a := range optPieces {
		bestPlays = append(bestPlays, simulatePlays(pl, a, 2)...)
	}

	bestPlay := selectPlay(pl, bestPlays)
	modutil.PrintSystem("executing play...")

	pieceCoord := PointToCoord(bestPlay.piece.Point)
	fmt.Printf("\nPlaying piece %s to %s\n\n", pieceCoord, PointToCoord(bestPlay.play.Dest))
	fmt.Print("Details:\n\n")
	fmt.Printf("- Piece: %s\n", pieceCoord)
	fmt.Printf("- Moves: %s\n", strings.Join(bestPlay.play.Breadcrumbs, ", "))

	if slays {
		coords := util.Map(bestPlay.play.Slays, func(i *Point) string { return PointToCoord(*i) })
		fmt.Printf("- Slays: %s\n", strings.Join(coords, ", "))
	}

	if bestPlay.play.IsKing && !bestPlay.piece.IsKing {
		fmt.Print("- Status: crowned\n")
	}

	////

	ExecutePlay(pl, bestPlay.piece, bestPlay.play)

	fmt.Println()
	modutil.PrintSystem("done\n")

	return Cont
}

func selectPlay(pl *player, pp []piecePlay) piecePlay {
	if len(pp) == 0 {
		panic("no possible plays found")
	}

	rand.Seed(time.Now().UnixNano())

	var maxScore *float64
	scoreMap := make(map[float64][]piecePlay)

	for _, a := range pp {
		if maxScore == nil || *maxScore < a.score {
			maxScore = &a.score
		}

		if _, ok := scoreMap[a.score]; ok {
			scoreMap[a.score] = append(scoreMap[a.score], a)
		} else {
			scoreMap[a.score] = []piecePlay{a}
		}
	}

	inx := rand.Intn(len(scoreMap[*maxScore]))

	return scoreMap[*maxScore][inx]
}

func simulatePlays(p *player, pi *Piece, depth int) []piecePlay {
	result := []piecePlay{}

	if depth == 0 {
		return result
	}

	// TODO: calculate score of leaving the piece without moving

	var maxScore *float64

	trees := CreateTreeMaps(*p, pi)
	plays := GetPiecePlays(trees)

	for x, a := range plays {
		score := 0.0

		if len(a.Slays) > 0 {
			for _, b := range a.Slays {
				ep, _, _ := p.board.GetPieceAt(*b)

				if ep.IsKing {
					score += float64(kingSlayMove)
				} else {
					score += float64(simpleSlayMove)
				}
			}
		}

		////

		if !pi.IsKing && a.IsKing {
			score += float64(crowningMove)
		} else if !pi.IsKing {
			crownDest := -1

			if p.Char == whiteChar {
				crownDest = 0
			} else {
				crownDest = boardSize - 1
			}

			score += float64(crowningMove) - math.Abs(float64(a.Dest.X-crownDest))
		} else {
			// TODO: find the closest piece and the best slaying spot
		}

		if depth-1 == 0 {
			if maxScore == nil || *maxScore < score {
				maxScore = &score
			}

			result = append(result, piecePlay{pi, a, score})
			continue
		}

		////

		player, piece, playsim := getStateCopy(*p, *pi)
		ExecutePlay(player, piece, playsim[x])

		var enemyOptPieces []*Piece

		if slayers := FilterSlayingOptions(*player.Enemy); len(slayers) > 0 {
			enemyOptPieces = slayers
		} else {
			enemyOptPieces = FilterSimpleOptions(*player.Enemy)
		}

		if len(enemyOptPieces) == 0 {
			score += float64(nullifyingMove)
		}

		enemyMoveCount := 0
		enemyScoreSum := 0.0

		for _, b := range enemyOptPieces {
			enemyResults := simulatePlays(player.Enemy, b, depth-1)

			for x := range enemyResults {
				enemyMoveCount++
				enemyScoreSum += enemyResults[x].score
			}
		}

		////

		if enemyMoveCount > 0 {
			score -= enemyScoreSum / float64(enemyMoveCount)
		}

		if maxScore == nil || *maxScore < score {
			maxScore = &score
		}

		result = append(result, piecePlay{pi, a, score})
	}

	return util.Filter(result, func(i piecePlay) bool { return i.score == *maxScore })
}

func getStateCopy(pl player, pi Piece) (*player, *Piece, []Play) {
	board := *pl.board

	board.white = make([]*Piece, len(pl.board.white))
	board.black = make([]*Piece, len(pl.board.black))

	enemy := *pl.Enemy
	pl.Enemy = &enemy
	pl.Enemy.Enemy = &pl

	for x := range pl.board.white {
		cp := *pl.board.white[x]
		board.white[x] = &cp
	}

	for x := range pl.board.black {
		cp := *pl.board.black[x]
		board.black[x] = &cp
	}

	pl.board = &board
	pl.Enemy.board = &board

	if pl.Char == whiteChar {
		pl.Pieces = &pl.board.white
		pl.Enemy.Pieces = &pl.board.black
	} else {
		pl.Pieces = &pl.board.black
		pl.Enemy.Pieces = &pl.board.white
	}

	piece, _, _ := pl.board.GetPieceAt(pi.Point)

	trees := CreateTreeMaps(pl, piece)
	plays := GetPiecePlays(trees)

	return &pl, piece, plays
}
