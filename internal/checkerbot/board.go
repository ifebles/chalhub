package checkerbot

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	whiteChar         = 'w'
	blackChar         = 'b'
	blankSpace        = '_'
	boardSize         = 8
	initialPieceCount = 12
)

type point struct {
	x, y int
}

type piece struct {
	point
	isKing bool
}

type board struct {
	initialized bool
	white       []piece
	black       []piece
}

var Board = &board{}

func (b *board) String() string {
	return b.Render()
}

func (b *board) Initialize() {
	b.initialized = true
	b.white = make([]piece, initialPieceCount)
	b.black = make([]piece, initialPieceCount)

	for x := 0; x < initialPieceCount; x++ {
		b.white[x] = piece{point{x/4 + 5, (x%4)*2 + (x/4)%2}, false}
		b.black[x] = piece{point{x / 4, (x%4)*2 + (1 - (x/4)%2)}, false}
	}
}

func (b *board) Render() string {
	boardRender := make([][]rune, boardSize)

	for x := range boardRender {
		boardRender[x] = []rune(strings.Repeat(string(blankSpace), boardSize))
	}

	populatePieces := func(pieces []piece, pieceChar rune, ch chan bool) {
		for x := range pieces {
			var char rune

			if pieces[x].isKing {
				char = unicode.ToUpper(pieceChar)
			} else {
				char = pieceChar
			}

			boardRender[pieces[x].x][pieces[x].y] = char
		}

		ch <- true
	}

	ch := make(chan bool, 2)

	go populatePieces(b.white, whiteChar, ch)
	go populatePieces(b.black, blackChar, ch)

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

func (b *board) Clear() {
	b.initialized = false
	b.white = nil
	b.black = nil
}
