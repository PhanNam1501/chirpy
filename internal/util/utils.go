package util

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func CleanProfanity(text string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range profaneWords {
		re := regexp.MustCompile(`(?i)\b` + word + `\b`)
		text = re.ReplaceAllString(text, "****")
	}
	return text
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(data)
}
