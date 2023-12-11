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

type apiConfig struct {
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

	
	

	apiCfg := apiConfig{
		DB: database.New(db),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	
	v1Router.Get("/users",  apiCfg.middlewareAuth(apiCfg.handlerGetUser) )
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Post("/feeds", apiCfg.handlerGetFeed)
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
