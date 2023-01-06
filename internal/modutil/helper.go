package modutil

import (
	"fmt"
	"math"
	"strings"

	"github.com/ifebles/chalhub/pkg/util"
)

func GetFormattedOptions(options []string, zeroText string, columns int) []string {
	if columns < 1 {
		panic(fmt.Sprintf("invalid column amount: %d", columns))
	}

	result := make([][]string, int(math.Ceil(float64(len(options))/float64(columns)))+1)

	for x := range options {
		if result[x/columns] == nil {
			result[x/columns] = make([]string, columns)
		}

		result[x/columns][x%columns] = fmt.Sprintf("\t%d) %s", x+1, options[x])
	}

	result[len(result)-1] = []string{fmt.Sprintf("\t0) %s", zeroText)}

	return util.Map(result, func(i []string) string {
		return strings.Join(i, "\t")
	})
}
