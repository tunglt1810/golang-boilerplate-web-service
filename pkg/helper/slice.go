package helper

func Filter[T any](source []T, filteredFunc func(item T) bool) (output []T) {
	for _, item := range source {
		if filteredFunc(item) {
			output = append(output, item)
		}
	}
	return output
}

func IndexOf[T comparable](source []T, item T) int {
	for i := 0; i < len(source); i++ {
		if item == source[i] {
			return i
		}
	}
	return -1
}
