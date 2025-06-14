package filter

import (
	"sort"
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

// Helper function to find the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper function to find the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Crops filters the given list of crops based on the query string.
// It scores each match using Jaro-Winkler similarity and returns the top 5 matches, sorted alphabetically.
func Crops(crops []string, query string) []string {
	if query == "" {
		return nil
	}
	type scoredCrop struct {
		crop  string
		score float64
	}
	scoredCrops := make([]scoredCrop, 0)
	for _, crop := range crops {
		score := jaroWinklerSimilarity(crop, query)
		scoredCrops = append(scoredCrops, scoredCrop{crop, score})
	}
	// Sort by score (descending) and then alphabetically
	sort.Slice(scoredCrops, func(i, j int) bool {
		if scoredCrops[i].score != scoredCrops[j].score {
			return scoredCrops[i].score > scoredCrops[j].score
		}
		return scoredCrops[i].crop < scoredCrops[j].crop
	})
	// Take the top 5 matches
	result := make([]string, 0)
	for i := 0; i < len(scoredCrops) && i < 5; i++ {
		result = append(result, scoredCrops[i].crop)
	}
	sort.Strings(result)
	return result
}
