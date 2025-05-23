package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"league-simulator/backend/db"
)

// UpdateMatchResult handles PUT /match/{id}.
// It allows manually editing the result of a match using updated scores.
func UpdateMatchResult(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for cross-origin access (e.g. frontend tools)
	setupCORS(w, r)

	// Only allow PUT method for updating
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract match ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/match/")
	matchID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	// Parse new scores from request body
	var update struct {
		HomeScore int `json:"home_score"`
		AwayScore int `json:"away_score"`
	}
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Determine the match result string
	result := "draw"
	if update.HomeScore > update.AwayScore {
		result = "win"
	} else if update.HomeScore < update.AwayScore {
		result = "loss"
	}

	// Update the match record in the database
	_, err = db.DB.Exec(`
		UPDATE matches
		SET home_score = ?, away_score = ?, result = ?
		WHERE id = ?
	`, update.HomeScore, update.AwayScore, result, matchID)

	if err != nil {
		http.Error(w, "Failed to update match", http.StatusInternalServerError)
		return
	}

	// Send confirmation response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Match updated successfully")
}
