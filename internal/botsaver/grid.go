package botsaver

import (
	"math/rand"
	"time"

	"chalhub/internal/modutil"
	"chalhub/pkg/util"
)

type Point struct {
	x, y int
}

func GetGridSize() int {
	for {
		gridSize, err := util.ReadInteger("Enter an odd number between 3 and 99: ")

		if err != nil {
			modutil.PrintAdvice("an integer was expected")
			continue
		}

		if gridSize < 3 || gridSize > 99 {
			modutil.PrintAdvice("an integer between 3 and 99 is needed; given: %d", gridSize)
			continue
		}

		if gridSize%2 == 0 {
			modutil.PrintAdvice("an odd integer is needed; given: %d (even)", gridSize)
			continue
		}

		return gridSize
	}
}

func GenerateGrid(size int) [][]string {
	middleValue := getMiddleValue(size)
	middlePoint, princessPoint := Point{middleValue, middleValue}, getPrincessPoint(size)

	result := make([][]string, size)

	for x := range result {
		result[x] = make([]string, size)

		for y := range result[x] {
			switch {
			case x == princessPoint.x && y == princessPoint.y:
				result[x][y] = "p"

			case x == middlePoint.x && y == middlePoint.y:
				result[x][y] = "m"

			default:
				result[x][y] = "-"
			}
		}
	}

	return result
}

func getMiddleValue(size int) int {
	return size / 2
}

func getPrincessPoint(size int) (result Point) {
	middleValue := getMiddleValue(size)
	rand.Seed(time.Now().UnixNano())

	for {
		result = Point{
			rand.Intn(size),
			rand.Intn(size),
		}

		if result.x != middleValue || result.y != middleValue {
			break
		}
	}

	return
}
