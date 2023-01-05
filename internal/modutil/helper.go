package modutil

import "fmt"

func GetFormattedOptions(options []string) []string {
	result := make([]string, 0, len(options)+1)

	for x := range options {
		result = append(result, fmt.Sprintf("\t%d) %s", x+1, options[x]))
	}

	result = append(result, "\t0) Exit")

	return result
}
