package checkerbot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type playerType string

const (
	Ai    playerType = "ai"
	Human playerType = "human"
)

type PlayMode uint8

const (
	PlayerVsAI PlayMode = iota
	PlayerVsPlayer
	AIvsAI
)

func (mode PlayMode) String() string {
	switch mode {
	case PlayerVsAI:
		return "Player vs AI"

	case PlayerVsPlayer:
		return "Player vs Player"

	case AIvsAI:
		return "AI vs AI"
	}

	return "unknown"
}

var PlayModes = [3]PlayMode{
	PlayerVsAI, PlayerVsPlayer, AIvsAI,
}

type player struct {
	Char  rune
	Enemy *player
	Type  playerType
}

var currentTurn = 1
var players = [2]player{}

func StartGame(mode PlayMode) player {
	if !Board.initialized {
		go Board.initialize()
	}

	switch mode {
	case PlayerVsAI:
		rand.Seed(time.Now().UnixNano())

		if rand.Intn(2) == 1 {
			players[0] = player{blackChar, &players[1], Human}
			players[1] = player{whiteChar, &players[0], Ai}
		} else {
			players[0] = player{blackChar, &players[1], Ai}
			players[1] = player{whiteChar, &players[0], Human}
		}

	case PlayerVsPlayer:
		players[0] = player{blackChar, &players[1], Human}
		players[1] = player{whiteChar, &players[0], Human}

	case AIvsAI:
		players[0] = player{blackChar, &players[1], Ai}
		players[1] = player{whiteChar, &players[0], Ai}

	default:
		panic("unknown mode")
	}

	return players[0]
}

func GetTurnNumber() int {
	return currentTurn
}

func EndTurn() player {
	currentTurn++

	return players[(currentTurn-1)%len(players)]
}

func GetPieces(pl player) []Piece {
	if pl.Char != whiteChar && pl.Char != blackChar {
		panic(fmt.Sprintf("unknown character: %s", string(pl.Char)))
	}

	var pieces []Piece

	if pl.Char == whiteChar {
		pieces = Board.white
	} else {
		pieces = Board.black
	}

	return pieces
}

/**func GetPiecePositions(pl player) []point {
	pieces := GetPieces(pl)
	points := util.Map(pieces, func(i Piece) point { return i.Point })

	return points
}

func GetPieceCoordinates(pl player) []string {
	pieces := GetPieces(pl)
	coords := util.Map(pieces, func(i Piece) string { return PointToCoord(i.Point) })

	return coords
}//*/

func PointToCoord(p point) string {
	if p.X < 0 || p.X >= boardSize || p.Y < 0 || p.Y >= boardSize {
		panic(fmt.Sprintf("invalid point: %v", p))
	}

	literal, numeral := rune('A'+p.Y), boardSize-p.X

	return fmt.Sprintf("%s%d", string(literal), numeral)
}

func CoordToPoint(c string) point {
	if len(c) != 2 {
		panic(fmt.Sprintf("invalid coordinate: %q", c))
	}

	literal := rune(strings.ToUpper(c)[0])

	if literal < 'A' || literal > 'A'+boardSize-1 {
		panic(fmt.Sprintf("invalid literal: %s", string(literal)))
	}

	numeral, err := strconv.Atoi(string(c[1]))

	if err != nil {
		panic(err)
	}

	if numeral <= 0 || numeral > boardSize {
		panic(fmt.Sprintf("invalid numeral: %d", numeral))
	}

	return point{boardSize - numeral, int(literal - 'A')}
}
