package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	internal "github.com/sakthi-lucia0567/rssagg/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *internal.Queries
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

	feed, err := urlToFeed("https://wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed)

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// conn, err := sql.Open("postgres", dbUrl)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	apiCfg := apiConfig{
		DB: internal.New(conn),
	}

	router := chi.NewRouter()

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness)
	v1Router.Get("/err", handleError)
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))

	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feed", apiCfg.handleGetFeeds)
	v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follow", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollow))
	v1Router.Delete("/feed_follow/{FeedFollowId}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))

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
