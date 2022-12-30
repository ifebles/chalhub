// Challenge inspiration from: https://www.hackerrank.com/challenges/saveprincess

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Point struct {
	x, y int
}

func main() {
	fmt.Print("** Welcome to the \"Save the Princess\" solver **\n\n")

	gridSize := getGridSize()

	fmt.Printf("> Generating %dx%d grid...\n", gridSize, gridSize)
	grid := generateGrid(gridSize)
	fmt.Print("> Done\n\n")

	pauseExecution()

	fmt.Println(grid)
}

func generateGrid(size int) [][]string {
	middleValue := getMiddleValue(size)
	middlePoint, princessPoint := Point{middleValue, middleValue}, getPrincessPoint(size)
	fmt.Println(princessPoint)

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

func readInteger(prompt string) (int, error) {
	var number int
	stdin := bufio.NewReader(os.Stdin)

	fmt.Print(prompt)

	_, err := fmt.Fscan(stdin, &number)
	stdin.Discard(stdin.Buffered())

	return number, err
}

func pauseExecution() {
	stdin := bufio.NewReader(os.Stdin)
	pauseMessage := "(Press 'Enter' to continue)"

	fmt.Print(pauseMessage)
	fmt.Fscanln(stdin)

	stdin.Discard(stdin.Buffered())
}

func getGridSize() int {
	for {
		gridSize, err := readInteger("Enter an odd number between 3 and 99: ")

		if err != nil {
			printAdvice("an integer was expected")
			continue
		}

		if gridSize < 3 || gridSize > 99 {
			printAdvice("an integer between 3 and 99 is needed; given: %d", gridSize)
			continue
		}

		if gridSize%2 == 0 {
			printAdvice("an odd integer is needed; given: %d (even)", gridSize)
			continue
		}

		return gridSize
	}
}

// printAdvice prints the given str to the stdout with an "Advice" tag.
// Accepts formatting.
func printAdvice(str string, args ...any) {
	fmt.Printf("Advice: "+str+"\n", args...)
}
