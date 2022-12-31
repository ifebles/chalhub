package botsaver

import (
	"fmt"
	"math"
	"strings"
)

const (
	up    = "UP"
	down  = "DOWN"
	right = "RIGHT"
	left  = "LEFT"
)

type bot struct {
	grid [][]rune
	path []string
}

var Bot = bot{}

func (b *bot) Feed(grid string) {
	rows := strings.Split(grid, "\n")
	b.grid = make([][]rune, len(rows))

	for x := range b.grid {
		b.grid[x] = []rune(rows[x])
	}
}

func (b *bot) DisplayPathToPrincess(cb func()) {
	b.calculatePath()
	cb()
	fmt.Println(strings.Join(b.path, "\n"))
}

func (b *bot) Clear() {
	b.grid = nil
	b.path = nil
}

func (b *bot) calculatePath() {
	if b.grid == nil {
		panic("no grid fed yet")
	}

	middlePoint, princessPoint := getPoints(b.grid)
	vertDist, horDist := princessPoint.x-middlePoint.x, princessPoint.y-middlePoint.y

	for x := 0.0; x < math.Abs(float64(vertDist)); x++ {
		if vertDist < 0 {
			b.path = append(b.path, up)
		} else {
			b.path = append(b.path, down)
		}
	}

	for x := 0.0; x < math.Abs(float64(horDist)); x++ {
		if horDist < 0 {
			b.path = append(b.path, left)
		} else {
			b.path = append(b.path, right)
		}
	}
}

func getPoints(grid [][]rune) (middle, princess Point) {
	middleValue := len(grid) / 2

	if grid[middleValue][middleValue] == middleChar {
		middle = Point{middleValue, middleValue}
	}

	for x := range grid {
		for y := range grid[x] {
			switch grid[x][y] {
			case princessChar:
				princess = Point{x, y}

				if middle != (Point{}) {
					return
				}

			case middleChar:
				if middle == (Point{}) {
					middle = Point{x, y}
				}

				if princess != (Point{}) {
					return
				}
			}
		}
	}

	return
}
