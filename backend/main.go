package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"

	"nklein.xyz/backend/routes"
	"nklein.xyz/backend/server"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL not set")
		os.Exit(1)
	}

	// Initialize Server (connects to DB)
	s, err := server.New(dbURL)
	if err != nil {
		panic(err)
	}

	// Router
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://atlas:2989",
			"https://nklein.xyz",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Mount grouped routes
	r.Mount("/api", routes.HealthRoutes(s))
	r.Mount("/api/landingcard", routes.LandingCardRoutes(s))

	addr := ":8080"
	fmt.Printf("🚀 Starting server on %s\n", addr)
	http.ListenAndServe(addr, r)
}
