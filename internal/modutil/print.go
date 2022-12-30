package modutil

import "fmt"

// PrintAdvice prints the given str to the stdout with an "Advice" tag.
// Accepts formatting.
func PrintAdvice(str string, args ...any) {
	fmt.Printf("Advice: "+str+"\n", args...)
}

// PrintSystem prints the given str to the stdout with a "System" tag.
// Accepts formatting.
func PrintSystem(str string, args ...any) {
	fmt.Printf("> "+str+"\n", args...)
}
