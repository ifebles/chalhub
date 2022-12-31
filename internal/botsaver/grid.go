package botsaver

import (
	"math/rand"
	"strings"
	"time"

	"github.com/ifebles/chalhub/internal/modutil"
	"github.com/ifebles/chalhub/pkg/util"
)

const (
	middleChar   = 'm'
	princessChar = 'p'
	blankSpace   = '-'
)

type Point struct {
	x, y int
}

type grid struct {
	size   int
	matrix [][]rune
}

var Grid = &grid{}

func (g *grid) String() string {
	rows := make([]string, len(g.matrix))

	for x := range g.matrix {
		rows[x] = string(g.matrix[x])
	}

	return strings.Join(rows, "\n")
}

func (g *grid) GetGridSize() int {
	for {
		var err error
		g.size, err = util.ReadInteger("Enter an odd number between 3 and 99: ")

		if err != nil {
			modutil.PrintAdvice("an integer was expected")
			continue
		}

		if g.size < 3 || g.size > 99 {
			modutil.PrintAdvice("an integer between 3 and 99 is needed; given: %d", g.size)
			continue
		}

		if g.size%2 == 0 {
			modutil.PrintAdvice("an odd integer is needed; given: %d (even)", g.size)
			continue
		}

		return g.size
	}
}

func (g *grid) GenerateGrid() bool {
	if g.size == 0 {
		return false
	}

	middleValue := g.getMiddleValue()
	middlePoint, princessPoint := Point{middleValue, middleValue}, g.getPrincessPoint()

	g.matrix = make([][]rune, g.size)

	for x := range g.matrix {
		g.matrix[x] = make([]rune, g.size)

		for y := range g.matrix[x] {
			switch {
			case x == princessPoint.x && y == princessPoint.y:
				g.matrix[x][y] = princessChar

			case x == middlePoint.x && y == middlePoint.y:
				g.matrix[x][y] = middleChar

			default:
				g.matrix[x][y] = blankSpace
			}
		}
	}

	return true
}

func (g *grid) Clear() {
	g.size = 0
	g.matrix = nil
}

func (g *grid) getMiddleValue() int {
	return g.size / 2
}

func (g *grid) getPrincessPoint() (result Point) {
	middleValue := g.getMiddleValue()
	rand.Seed(time.Now().UnixNano())

	for {
		result = Point{
			rand.Intn(g.size),
			rand.Intn(g.size),
		}

		if result.x != middleValue || result.y != middleValue {
			break
		}
	}

	return
}
