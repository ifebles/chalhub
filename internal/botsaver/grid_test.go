package botsaver

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestSizeInput(t *testing.T) {
	tGrid := *Grid

	cases := []struct {
		name   string
		inputs []string
		want   int
	}{
		{"get grid size", []string{"99"}, 99},
		{"fail after string input", []string{"asdf", "qwerty", "zxcv", ",./"}, -1},
		{"fail after float input", []string{"3.7", "3.0", "4.2", "6.5"}, -1},
		{"fail after low input", []string{"-1", "0", "1", "2"}, -1},
		{"fail after high input", []string{"101", "1001", "10001"}, -1},
		{"fail after even input", []string{"4", "80", "14", "98"}, -1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tmpStdinFile, err := createTempFile("testSizeInput")

			if err != nil {
				log.Fatal(err)
			}

			defer os.Remove(tmpStdinFile.Name())

			tmpStdoutFile, err := createTempFile("testSizeOutput")

			if err != nil {
				log.Fatal(err)
			}

			defer os.Remove(tmpStdoutFile.Name())

			////

			oldStdin, oldStdout := os.Stdin, os.Stdout

			defer func() { os.Stdin = oldStdin }()

			os.Stdin = tmpStdinFile

			////

			var result int

			for _, a := range c.inputs {
				if err := writeToFile(tmpStdinFile, a); err != nil {
					log.Fatal(err)
				}

				func() {
					os.Stdout = tmpStdoutFile
					defer func() { os.Stdout = oldStdout }()

					result = tGrid.GetGridSize(1)
				}()
			}

			if want, got := c.want, result; want != got {
				t.Errorf("unexpected value from stdin; want: %d, got: %d", want, got)
			}
		})
	}
}

func TestInvalidCallToGetSize(t *testing.T) {
	tGrid := *Grid

	defer func() {
		if err := recover(); err == nil {
			t.Error("expected a panic to occur")
		}
	}()

	tGrid.GetGridSize(0)
}

func TestGridGeneration(t *testing.T) {
	const size = 3
	tGrid := *Grid

	result := tGrid.GenerateGrid()

	if want, got := false, result; want != got {
		t.Errorf("unexpected result from `GenerateGrid` method; want: %v, got: %v", want, got)
	}

	tGrid.size = size

	if tGrid.matrix != nil {
		t.Errorf("unexpected value for the matrix; want: %v, got: %v", nil, tGrid.matrix)
	}

	result = tGrid.GenerateGrid()

	if want, got := true, result; want != got {
		t.Fatalf("unexpected result from `GenerateGrid` method; want: %v, got: %v", want, got)
	}

	if want, got := size, len(tGrid.matrix); want != got {
		t.Errorf("unexpected matrix size; want: %d, got: %d", want, got)
	}

	if want, got := size, len(tGrid.matrix[0]); want != got {
		t.Errorf("unexpected array size for matrix; want: %d, got: %d", want, got)
	}

	////

	foundBot, foundPrincess := false, false

	for x := range tGrid.matrix {
		for y := range tGrid.matrix[x] {
			switch tGrid.matrix[x][y] {
			case princessChar:
				foundPrincess = true

				if foundBot {
					break
				}
			case middleChar:
				foundBot = true

				if foundPrincess {
					break
				}
			}
		}

		if foundBot && foundPrincess {
			break
		}
	}

	////

	if !foundBot {
		t.Errorf("missing required bot character (%v) in the grid", middleChar)
	}

	if !foundPrincess {
		t.Errorf("missing required princess character (%v) in the grid", princessChar)
	}
}

func TestGridString(t *testing.T) {
	const size = 3
	tGrid := *Grid

	tGrid.size = size
	tGrid.GenerateGrid()

	result := fmt.Sprint(&tGrid)

	if result == "" {
		t.Error("invalid string version of the grid")
	}

	if lineBreaks := strings.Count(result, "\n"); lineBreaks != size-1 {
		t.Errorf(
			"unexpected amount of line breaks from the grid's resulting string; want: %d, got: %d",
			size-1,
			lineBreaks,
		)
	}

	if want, got := size*size+(size-1), len(result); want != got {
		t.Errorf("unexpected string length; want: %d, got: %d", want, got)
	}

	if !strings.ContainsRune(result, middleChar) {
		t.Errorf("missing bot character (%v) from the resulting string", middleChar)
	}

	if strings.Count(result, string(middleChar)) > 1 {
		t.Errorf("more than one occurrence for the bot character (%v) in the resulting string", middleChar)
	}

	if !strings.ContainsRune(result, princessChar) {
		t.Errorf("missing bot character (%v) from the resulting string", princessChar)
	}

	if strings.Count(result, string(princessChar)) > 1 {
		t.Errorf("more than one occurrence for the bot character (%v) in the resulting string", princessChar)
	}

	if want, got := len(result)-(size-1)-2, strings.Count(result, string(blankSpace)); want != got {
		t.Errorf("unexpected amount of blank spaces found; want: %d, got: %d", want, got)
	}
}

func TestGridClear(t *testing.T) {
	tGrid := *Grid

	tGrid.size = 3
	tGrid.GenerateGrid()
	tGrid.Clear()

	if want, got := 0, tGrid.size; want != got {
		t.Errorf("unexpected size; want: %d, got: %d", want, got)
	}

	if tGrid.matrix != nil {
		t.Errorf("unexpected matrix value; want: %v, got: %v", nil, tGrid.matrix)
	}
}

func createTempFile(name string) (*os.File, error) {
	file, err := os.CreateTemp("", name)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func writeToFile(file *os.File, content string) error {
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Write([]byte(content)); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}
