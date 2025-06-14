package filter

import (
	"strings"
)

// jaroWinklerSimilarity calculates the Jaro-Winkler similarity between two strings.
// Returns a score between 0 and 1, where 1 is a perfect match.
func jaroWinklerSimilarity(s1, s2 string) float64 {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)
	if s1 == s2 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Jaro distance calculation
	matchDistance := (max(len(s1), len(s2)) / 2) - 1
	s1Matches := make([]bool, len(s1))
	s2Matches := make([]bool, len(s2))
	matches := 0
	transpositions := 0

	for i := 0; i < len(s1); i++ {
		start := max(0, i-matchDistance)
		end := min(len(s2), i+matchDistance+1)
		for j := start; j < end; j++ {
			if !s2Matches[j] && s1[i] == s2[j] {
				s1Matches[i] = true
				s2Matches[j] = true
				matches++
				break
			}
		}
	}

	if matches == 0 {
		return 0.0
	}

	// Count transpositions
	k := 0
	for i := 0; i < len(s1); i++ {
		if s1Matches[i] {
			for !s2Matches[k] {
				k++
			}
			if s1[i] != s2[k] {
				transpositions++
			}
			k++
		}
	}
	transpositions /= 2

	// Calculate Jaro distance
	jaroDistance := (float64(matches)/float64(len(s1)) + float64(matches)/float64(len(s2)) + float64(matches-transpositions)/float64(matches)) / 3.0

	// Jaro-Winkler adjustment
	prefixLength := 0
	for i := 0; i < min(4, min(len(s1), len(s2))); i++ {
		if s1[i] == s2[i] {
			prefixLength++
		} else {
			break
		}
	}

	jaroWinklerDistance := jaroDistance + float64(prefixLength)*0.1*(1.0-jaroDistance)
	return jaroWinklerDistance
}
