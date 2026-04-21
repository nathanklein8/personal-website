package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// Request logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			log.Printf("[%s] %s %s", req.Method, req.URL.Path, req.RemoteAddr)
			next.ServeHTTP(w, req)
			log.Printf("Completed %s %s in %v", req.Method, req.URL.Path, time.Since(start))
		})
	})

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	// Mount route sub-routers
	r.Mount("/api/health", routes.HealthRoutes(s))
	r.Mount("/api/landingcard", routes.LandingCardRoutes(s))
	r.Mount("/api/projects", routes.ProjectRoutes(s))
	r.Mount("/api/photos", routes.PhotoRoutes(s))

	addr := ":8080"
	fmt.Printf("🚀 Starting server on %s\n", addr)
	http.ListenAndServe(addr, r)
}
