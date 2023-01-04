package botsaver

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestFeed(t *testing.T) {
	tBot, tGrid := Bot, *Grid

	if tBot.grid != nil {
		t.Errorf("unexpected value for grid; want: %v, got: %v", nil, tBot.grid)
	}

	if tBot.path != nil {
		t.Errorf("unexpected value for path; want: %v, got: %v", nil, tBot.path)
	}

	tGrid.size = 3
	tGrid.GenerateGrid()

	tBot.Feed(fmt.Sprint(&tGrid))

	if want, got := fmt.Sprint(tGrid.matrix), fmt.Sprint(tBot.grid); want != got {
		t.Errorf("unexpected value for grid; want: %v, got: %v", want, got)
	}
}

func TestBotClear(t *testing.T) {
	tBot, tGrid := Bot, *Grid

	tGrid.size = 3
	tGrid.GenerateGrid()

	tBot.Feed(fmt.Sprint(&tGrid))

	////

	func() {
		tmpFile, err := createTempFile("botClear")

		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpFile.Name())

		oldStdout := os.Stdout

		func() {
			os.Stdout = tmpFile
			defer func() { os.Stdout = oldStdout }()

			tBot.DisplayPathToPrincess(func() {})
		}()
	}()

	////

	if tBot.grid == nil {
		t.Error("unexpected nil value in grid")
	}

	if tBot.path == nil {
		t.Error("unexpected nil value in path")
	}

	if len(tBot.path) == 0 {
		t.Error("unexpected empty value for path")
	}

	tBot.Clear()

	if tBot.grid != nil {
		t.Errorf("unexpected value for grid; want: %v, got: %v", nil, tBot.grid)
	}

	if tBot.path != nil {
		t.Errorf("unexpected value for path; want: %v, got: %v", nil, tBot.path)
	}
}

func TestDisplayPathToPrincess(t *testing.T) {
	tBot, tGrid := Bot, *Grid

	tGrid.size = 3
	middleValue := tGrid.getMiddleValue()

	cases := []struct {
		name          string
		princessPoint Point
		want          string
	}{
		{"path to upper-left corner", Point{0, 0}, fmt.Sprintf("%s\n%s\n", up, left)},
		{"path to upper-middle side", Point{0, middleValue}, fmt.Sprintf("%s\n", up)},
		{"path to upper-right corner", Point{0, tGrid.size - 1}, fmt.Sprintf("%s\n%s\n", up, right)},
		{"path to left-middle side", Point{middleValue, 0}, fmt.Sprintf("%s\n", left)},
		{"path to right-middle side", Point{middleValue, tGrid.size - 1}, fmt.Sprintf("%s\n", right)},
		{"path to lower-left corner", Point{tGrid.size - 1, 0}, fmt.Sprintf("%s\n%s\n", down, left)},
		{"path to lower-middle side", Point{tGrid.size - 1, middleValue}, fmt.Sprintf("%s\n", down)},
		{"path to lower-right corner", Point{tGrid.size - 1, tGrid.size - 1}, fmt.Sprintf("%s\n%s\n", down, right)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tGrid.GenerateGrid()
			setPrincessPoint(&tGrid, &c.princessPoint)

			tBot.Feed(fmt.Sprint(&tGrid))

			////

			tmpFile, err := createTempFile("displayPathToPrincess")

			if err != nil {
				log.Fatalf("error creating tmpfile: %s", err)
			}

			defer os.Remove(tmpFile.Name())

			////

			oldStdout := os.Stdout

			func() {
				os.Stdout = tmpFile
				defer func() { os.Stdout = oldStdout }()

				tBot.DisplayPathToPrincess(func() {})
			}()

			if tBot.path == nil || len(tBot.path) == 0 {
				t.Errorf("expected an array value with length > 0 for path, got: %v", tBot.path)
			}

			////

			result, err := readFile(tmpFile)

			if err != nil {
				log.Fatalf("error reading tmpfile: %s", err)
			}

			////

			if want, got := c.want, result; want != got {
				t.Errorf("unexpected path result; want: %q, got: %q", want, got)
			}

			tBot.Clear()
		})
	}
}

func TestInvalidDiplayPathToPrincess(t *testing.T) {
	tBot := Bot

	defer func() {
		if err := recover(); err == nil {
			t.Error("expected a panic to occur")
		}
	}()

	tBot.DisplayPathToPrincess(func() {})
}

func setPrincessPoint(grid *grid, point *Point) {
	cleared := false

	for x := range grid.matrix {
		for y := range grid.matrix[x] {
			if grid.matrix[x][y] == princessChar {
				grid.matrix[x][y] = blankSpace
				cleared = true

				break
			}
		}

		if cleared {
			break
		}
	}

	grid.matrix[point.x][point.y] = princessChar
}

func readFile(file *os.File) (string, error) {
	var size int64

	if stat, err := file.Stat(); err == nil {
		size = stat.Size()
	} else {
		return "", err
	}

	content := make([]byte, size)

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	if _, err := file.Read(content); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", content), nil
}
