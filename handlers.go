package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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

func (apiCfg *apiConfig)handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct{
		Name string `json:"name"`
		Url string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	body := params{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
      ID: uuid.New(),
	  CreatedAt: time.Now().UTC(),
	  UpdatedAt: time.Now().UTC(),
	  Name:      body.Name,
	  Url: body.Url,
	  UserID: user.ID,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Could not create user: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig)handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct{
		FeedID uuid.UUID `json:"feed_id"`
		
	}

	decoder := json.NewDecoder(r.Body)
	body := params{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
      ID: uuid.New(),
	  CreatedAt: time.Now().UTC(),
	  UpdatedAt: time.Now().UTC(),
	 
	  UserID: user.ID,
	  FeedID: body.FeedID,
	})

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Could not create feed: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig)handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Could not get feed follows: %v", err))
		return
	}
	respondWithJSON(w, 200, feedFollows)
}

func (apiCfg *apiConfig)handlerGetFeed(w http.ResponseWriter, r *http.Request) {
	

	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Could not get feeds: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedToFeeds(feeds))
}
func (apiCfg *apiConfig)handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	

	respondWithJSON(w, 200, databaseUserToUser(user))

}


func (cfg *apiConfig) handlerFeedFollowDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	err = cfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

