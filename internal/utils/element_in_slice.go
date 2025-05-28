package utils

func ElementInSlice[T comparable](s []T, element T) bool {
	for _, v := range s {
		if v == element {
			return true
		}
	}
	return false
}
