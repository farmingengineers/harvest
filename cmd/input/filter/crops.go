package filter

import (
	"sort"
)

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
