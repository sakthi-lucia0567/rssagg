package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port is not fount in the environment")
	}

	router := chi.NewRouter()

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness)
	v1Router.Get("/err", handleError)

	router.Mount("/v1", v1Router)

	log.Printf("Server starting on port %v", port)
	error := server.ListenAndServe()
	if err != nil {
		log.Fatal(error)
	}

	router.Use(middleware.RouteHeaders().
		Route("Origin", "https://app.skyweaver.net", cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://api.skyweaver.net", "http://*", "https://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			AllowCredentials: false, // <----------<<< allow credentials
			MaxAge:           300,
		})).
		Route("Origin", "*", cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type"},
			AllowCredentials: false, // <----------<<< do not allow credentials
		})).
		Handler)
}
