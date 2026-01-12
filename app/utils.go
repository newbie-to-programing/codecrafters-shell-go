package main

func NotIn[T comparable](slice1, slice2 []T) []T {
	set1 := make(map[T]struct{})
	for _, item := range slice1 {
		set1[item] = struct{}{}
	}

	var notInElements []T
	seenInResult := make(map[T]bool)

	for _, item := range slice2 {
		if _, found := set1[item]; !found {
			if !seenInResult[item] {
				notInElements = append(notInElements, item)
				seenInResult[item] = true
			}
		}
	}

	return notInElements
}
