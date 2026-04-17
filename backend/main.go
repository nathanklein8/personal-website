package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
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

	// Mount route sub-routers
	r.Mount("/api", routes.HealthRoutes(s))
	r.Mount("/api", routes.LandingCardRoutes(s))
	r.Mount("/api", routes.ProjectRoutes(s))
	r.Mount("/api", routes.PhotoRoutes(s))

	addr := ":8080"
	fmt.Printf("🚀 Starting server on %s\n", addr)
	http.ListenAndServe(addr, r)
}
