package modutil

import "github.com/ifebles/chalhub/pkg/util"

func GetIntOptionFromUser(attemptLimit, valueLimit int) int {
	if attemptLimit <= 0 {
		panic("invalid limit given")
	}

	for x := 0; x < attemptLimit; x++ {
		result, err := util.ReadInteger(">> ")

		if err != nil {
			PrintAdvice("an integer was expected")
			continue
		}

		if min, max := 0, valueLimit; result < min || result > max {
			PrintAdvice("an integer between %d and %d was expected", min, max)
			continue
		}

		return result
	}

	return 0
}
