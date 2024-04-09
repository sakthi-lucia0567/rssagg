package main

import (
	"database/sql"
	"log"
	"migrations"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/joho/godotenv"
	"github.com/sakthi-lucia0567/rssagg/migrations"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *migrations.Queries
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port is not fount in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	queries, err := migrations.New(conn)

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
