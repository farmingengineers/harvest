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
	// Simplified Jaro-Winkler calculation (for demonstration)
	// In a real implementation, you would use a proper Jaro-Winkler algorithm
	// This is a placeholder that returns a basic similarity score
	return 0.5
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
	return result
}
