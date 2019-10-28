package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shortlink/utils"

	"github.com/gorilla/mux"

	"gopkg.in/go-playground/validator.v9"
)

//ShortRequest for request: url, expiration_in_minutes
type ShortRequest struct {
	URL                 string `json:"url" validate:"required"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

//ShortlinkRequest :shortlink
type ShortlinkRequest struct {
	Shortlink string `json:"shortlink"`
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req ShortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, utils.StatusError{Code: http.StatusBadRequest,
			Err: fmt.Errorf("json decode err %v", r.Body)})
		return
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		respondWithError(w, utils.StatusError{Code: http.StatusBadRequest,
			Err: fmt.Errorf("json validate err %v", req)})
		return
	}

	defer r.Body.Close()
	s, err := a.Config.S.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusCreated, ShortlinkRequest{Shortlink: s})
	}

}

func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	s := url.Get("shortlink")

	d, err := a.Config.S.ShortLinkInfo(s)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusOK, d)
	}

}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u, err := a.Config.S.Unshorten(vars["shortlink"])
	if err != nil {
		respondWithError(w, err)
	} else {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		return
	}
}
