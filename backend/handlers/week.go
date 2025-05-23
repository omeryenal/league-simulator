package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"league-simulator/backend/db"
)

const MaxWeek = 12

// GetCurrentWeek handles GET /week/current.
// It returns the next week number that should be played in the league.
// If no matches have been played yet, it starts from week 4.
func GetCurrentWeek(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow frontend interaction
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Query the last played week
	row := db.DB.QueryRow(`SELECT MAX(week) FROM matches`)
	var maxWeek sql.NullInt64
	err := row.Scan(&maxWeek)
	if err != nil {
		http.Error(w, "Failed to get current week", http.StatusInternalServerError)
		return
	}

	var week int
	if !maxWeek.Valid || maxWeek.Int64 < 4 {
		// If no matches or only up to week 3, start at week 4
		week = 4
	} else if maxWeek.Int64 >= MaxWeek {
		// If the season is complete, return a sentinel value
		week = MaxWeek + 1
	} else {
		// Otherwise, return the next week to be played
		week = int(maxWeek.Int64) + 1
	}

	// Return the week value as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"week": week})
}
