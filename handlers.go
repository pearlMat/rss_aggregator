package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
    "github.com/pearlMat/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 400, "something went wrong")
}

func (apiCfg *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct{
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	body := params{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
      ID: uuid.New(),
	  CreatedAt: time.Now().UTC(),
	  UpdatedAt: time.Now().UTC(),
	  Name:      body.Name,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Could not create user: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseUserToUser(user))
}