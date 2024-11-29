package slices

func HasDuplicates[T comparable](slice []T) bool {
	mp := make(map[T]struct{}, len(slice))
	for _, v := range slice {
		if _, ok := mp[v]; ok {
			return true
		}
		mp[v] = struct{}{}
	}

	return false
}
