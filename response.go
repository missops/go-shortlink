package main

import (
	"encoding/json"
	"log"
	"net/http"
	"shortlink/utils"
)

func respondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case utils.Error:
		log.Printf("HTTP %d - %s ", e.Status(), e)
		respondWithJSON(w, e.Status(), e.Error())
	default:
		respondWithJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
