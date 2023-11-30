package main

import (
  
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main(){
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
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
	  router.Mount("v1", v1Router)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
 
	err := srv.ListenAndServe()
	if err != nil{
		log.Fatal(err)
	}
    log.Print("server starting on port %v", port)

	
}