package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pearlMat/rss_aggregator/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}
func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == ""{
		log.Fatal("Port not found!")
	}
	dbUrl := os.Getenv("CONN")
	if dbUrl == ""{
		log.Fatal("dbUrl not found!")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil{
		log.Fatal("Cant connect to database", err)
	}

	
	

	apiCfg := ApiConfig{
		DB: database.New(db),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	
	v1Router.Get("/users",  apiCfg.middlewareAuth(apiCfg.handlerGetUser) )
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Post("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))
	router.Mount("v1", v1Router)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("server starting on port %v", port)

}
