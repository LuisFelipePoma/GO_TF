package utils

import ()

func SplitRanges(totalMovies, numParts int) [][2]int {
	chunkSize := (totalMovies + numParts - 1) / numParts
	var ranges [][2]int
	for i := 0; i < totalMovies; i += chunkSize {
		end := i + chunkSize
		if end > totalMovies {
			end = totalMovies
		}
		ranges = append(ranges, [2]int{i, end})
	}
	return ranges
}
