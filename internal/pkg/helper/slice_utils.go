package helper

func SliceDiff[T comparable](a, b []T) []T {
	var result []T
	for _, bItem := range b {
		ok := true
		for _, aItem := range a {
			if bItem == aItem {
				ok = false
			}
		}
		if ok {
			result = append(result, bItem)
		}
	}

	return result
}
