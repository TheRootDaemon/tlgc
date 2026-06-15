package slice

// Dedup returns a new slice with duplicate elements removed,
// preserving the order of first occurrence.
func Dedup[T comparable](s []T) []T {
	if len(s) == 0 {
		return nil
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(s))

	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
