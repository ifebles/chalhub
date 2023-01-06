package util

func Map[T1 any, T2 any](collection []T1, transform func(item T1) T2) []T2 {
	result := make([]T2, 0, len(collection))

	for _, a := range collection {
		result = append(result, transform(a))
	}

	return result
}

func Find[T any](collection []T, finder func(item T) bool) (*T, bool) {
	for _, a := range collection {
		if finder(a) {
			return &a, true
		}
	}

	return nil, false
}
