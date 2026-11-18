package main

import "strings"

func formatPlace(raw string) string {
	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		return raw
	}

	city := titleCase(strings.ReplaceAll(parts[0], "_", " "))

	country := strings.ReplaceAll(parts[1], "_", " ")
	if len(country) <= 3 {
		country = strings.ToUpper(country)
	} else {
		country = titleCase(country)
	}

	return city + ", " + country
}

func titleCase(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}

func filterArtists(artists []Artist, query string) []Artist {
	query = strings.ToLower(strings.TrimSpace(query))

	var result []Artist
	for _, a := range artists {
		if strings.Contains(strings.ToLower(a.Name), query) {
			result = append(result, a)
		}
	}
	return result
}
