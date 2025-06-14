package filter

import "strings"

// Crops filters the given list of crops based on the query string.
// It returns the top 10 matches (or fewer if there aren't that many).
func Crops(crops []string, query string) []string {
	if query == "" {
		return nil
	}
	filtered := make([]string, 0)
	for _, crop := range crops {
		if strings.Contains(strings.ToLower(crop), strings.ToLower(query)) {
			filtered = append(filtered, crop)
			if len(filtered) >= 10 {
				break
			}
		}
	}
	return filtered
}
