package booleans

// All will return false, if any result is false, otherwise return true
func All(result ...bool) bool {
	for _, r := range result {
		if !r {
			return false
		}
	}

	return true
}

// Any will return true if any result is true, otherwise return false
func Any(result ...bool) bool {
	for _, r := range result {
		if r {
			return true
		}
	}

	return false
}
