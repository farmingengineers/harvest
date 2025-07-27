package filter

import (
	"sort"
	"strings"
)

// Crops filters the given list of crops based on the query string.
// It scores each match using Jaro-Winkler similarity and returns the top 5 matches, sorted alphabetically.
func Crops(crops []string, query string, limit int) []string {
	if query == "" {
		return nil
	}
	type scoredCrop struct {
		crop  string
		score float64
	}
	queryParts := splitCrop(query)
	scoredCrops := make([]scoredCrop, 0)
	for _, crop := range crops {
		totalScore := 0.0
		for _, queryWord := range queryParts {
			bestScore := 0.0
			for _, cropWord := range splitCrop(crop) {
				score := jaroWinklerSimilarity(cropWord, queryWord)
				bestScore = max(score, bestScore)
			}
			totalScore += bestScore
		}
		scoredCrops = append(scoredCrops, scoredCrop{crop, totalScore})
	}
	// Sort by score (descending) and then alphabetically
	sort.Slice(scoredCrops, func(i, j int) bool {
		if scoredCrops[i].score != scoredCrops[j].score {
			return scoredCrops[i].score > scoredCrops[j].score
		}
		return scoredCrops[i].crop < scoredCrops[j].crop
	})
	// Take the top 5 matches
	result := make([]string, 0, limit)
	for i := 0; i < len(scoredCrops) && i < limit; i++ {
		result = append(result, scoredCrops[i].crop)
	}
	return result
}

func splitCrop(crop string) []string {
	if crop == "" {
		return nil
	}
	res := make([]string, 0, strings.Count(crop, ",")+strings.Count(crop, " ")+1)
	nextComma := strings.Index(crop, ",")
	nextSpace := strings.Index(crop, " ")
	for nextComma != -1 || nextSpace != -1 {
		if nextComma == -1 || (nextSpace != -1 && nextSpace < nextComma) {
			res = append(res, crop[:nextSpace])
			crop = crop[nextSpace+1:]
		} else {
			res = append(res, crop[:nextComma])
			crop = crop[nextComma+1:]
		}
		for crop != "" && crop[0] == ' ' {
			crop = crop[1:]
		}
		nextComma = strings.Index(crop, ",")
		nextSpace = strings.Index(crop, " ")
	}
	if crop != "" {
		res = append(res, crop)
	}
	return res
}
