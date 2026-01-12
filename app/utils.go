package main

// Intersection returns the intersection of two slices.
func Intersection[T comparable](slice1, slice2 []T) []T {
	// Create a map to store elements of the first slice for quick lookups.
	// Use struct{}{} as a value for minimum memory usage (zero-sized type).
	set1 := make(map[T]struct{})
	for _, item := range slice1 {
		set1[item] = struct{}{}
	}

	// Create a slice for the common elements and another map to track
	// elements already added to the result to avoid duplicates in the output.
	var commonElements []T
	seenInResult := make(map[T]bool)

	// Iterate over the second slice and check if each element exists in set1.
	for _, item := range slice2 {
		if _, found := set1[item]; found {
			// If found and not already added to the result slice, add it.
			if !seenInResult[item] {
				commonElements = append(commonElements, item)
				seenInResult[item] = true
			}
		}
	}

	return commonElements
}
