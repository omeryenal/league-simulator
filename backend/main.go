package main

import (
	"fmt"
	"log"
	"net/http"

	"league-simulator/backend/db"
	"league-simulator/backend/handlers"
)

// withCORS is a simple middleware to enable CORS headers.
// Useful during local development when frontend and backend run on different origins.
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

func main() {
	// Initialize the database and create tables if they don't exist
	db.InitDB()
	fmt.Println("✅ Database connected and tables initialized.")

	// Health check endpoint
	http.HandleFunc("/ping", withCORS(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	}))

	// League-related endpoints
	http.HandleFunc("/teams", withCORS(handlers.GetTeams))                 // GET
	http.HandleFunc("/standings", withCORS(handlers.GetStandings))         // GET
	http.HandleFunc("/week/current", withCORS(handlers.GetCurrentWeek))    // GET
	http.HandleFunc("/simulate/next", withCORS(handlers.SimulateNextWeek)) // POST
	http.HandleFunc("/simulate/all", withCORS(handlers.SimulateAll))       // POST
	http.HandleFunc("/reset", withCORS(handlers.ResetSeason))              // POST
	http.HandleFunc("/results/week/", withCORS(handlers.GetWeekResults))   // GET
	http.HandleFunc("/predictions", withCORS(handlers.GetPredictions))     // GET

	// Manual match control
	http.HandleFunc("/match", withCORS(handlers.CreateMatch))        // POST /match
	http.HandleFunc("/match/", withCORS(handlers.UpdateMatchResult)) // PUT /match/{id}

	// Start the server
	log.Println("🚀 Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
