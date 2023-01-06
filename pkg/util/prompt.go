package util

import (
	"bufio"
	"fmt"
	"os"
)

func ReadInteger(prompt string) (int, error) {
	var number int
	stdin := bufio.NewReader(os.Stdin)

	fmt.Print(prompt)

	_, err := fmt.Fscan(stdin, &number)
	stdin.Discard(stdin.Buffered())

	return number, err
}

func ReadString(prompt string) (string, error) {
	var response string
	stdin := bufio.NewReader(os.Stdin)

	fmt.Print(prompt)

	_, err := fmt.Fscan(stdin, &response)
	stdin.Discard(stdin.Buffered())

	return response, err
}

func PauseExecution() {
	stdin := bufio.NewReader(os.Stdin)
	pauseMessage := "(Press 'Enter' to continue)"

	fmt.Print(pauseMessage)
	fmt.Fscanln(stdin)

	stdin.Discard(stdin.Buffered())
}
