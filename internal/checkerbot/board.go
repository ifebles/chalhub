package checkerbot

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/ifebles/chalhub/pkg/util"
)

const (
	whiteChar         = 'w'
	blackChar         = 'b'
	blankSpace        = '_'
	boardSize         = 8
	initialPieceCount = 12
)

type point struct {
	X, Y int
}

type Piece struct {
	Point  point
	IsKing bool
}

type board struct {
	initialized bool
	white       []Piece
	black       []Piece
}

var Board = &board{}

func (b *board) String() string {
	return b.Render()
}

func (b *board) Render() string {
	boardRender := make([][]rune, boardSize)
	ch := make(chan bool, 2)

	for x := range boardRender {
		boardRender[x] = []rune(strings.Repeat(string(blankSpace), boardSize))
	}

	populatePieces := func(pieces []Piece, pieceChar rune) {
		for x := range pieces {
			var char rune

			if pieces[x].IsKing {
				char = unicode.ToUpper(pieceChar)
			} else {
				char = pieceChar
			}

			boardRender[pieces[x].Point.X][pieces[x].Point.Y] = char
		}

		ch <- true
	}

	go populatePieces(b.white, whiteChar)
	go populatePieces(b.black, blackChar)

	////

	rows := make([]string, len(boardRender)+1)

	letters := ""

	for x := 'A'; x < 'A'+boardSize; x++ {
		letters += string(x)
	}

	<-ch
	<-ch

	for x := range boardRender {
		rows[x] = fmt.Sprintf("%d | %s", boardSize-x, strings.Join(strings.Split(string(boardRender[x]), ""), " "))
	}

	rows[len(rows)-1] = fmt.Sprintf("    %s", strings.Join(strings.Split(letters, ""), " "))

	////

	return strings.Join(rows, "\n")
}

func (b *board) GetPieceAt(p point) (*Piece, bool, error) {
	if p.X < 0 || p.Y < 0 || p.X >= boardSize || p.Y >= boardSize {
		return nil, false, fmt.Errorf("out of bounds")
	}

	ch := make(chan *Piece, 2)

	search := func(col []Piece) {
		found, _ := util.FindPtr(col, func(i Piece) bool {
			return i.Point.X == p.X && i.Point.Y == p.Y
		})

		ch <- found
	}

	go search(b.white)
	go search(b.black)

	result := <-ch

	if result != nil {
		return result, true, nil
	}

	result = <-ch

	return result, result != nil, nil
}

func (b *board) Clear() {
	b.initialized = false
	b.white = nil
	b.black = nil
}

func (b *board) initialize() {
	b.initialized = true
	b.white = make([]Piece, initialPieceCount)
	b.black = make([]Piece, initialPieceCount)

	for x := 0; x < initialPieceCount; x++ {
		b.white[x] = Piece{point{x/4 + 5, (x%4)*2 + (x/4)%2}, false}
		b.black[x] = Piece{point{x / 4, (x%4)*2 + (1 - (x/4)%2)}, false}
	}
}
