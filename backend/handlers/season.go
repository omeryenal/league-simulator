package handlers

import (
	"net/http"

	"league-simulator/backend/db"
)

// ResetSeason handles POST /reset
// It deletes all match results after week 4, effectively restarting the season from week 5.
func ResetSeason(w http.ResponseWriter, r *http.Request) {
	// Allow cross-origin requests (e.g. from frontend or Postman)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST method for this endpoint
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Delete all matches after week 4 to reset the league state
	_, err := db.DB.Exec("DELETE FROM matches WHERE week > ?", 4)
	if err != nil {
		http.Error(w, "Failed to reset season: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Season reset successful", "week": 5}`))
}
