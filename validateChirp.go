package main

import (
	"strings"
)

func replaceProfaneWords(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
