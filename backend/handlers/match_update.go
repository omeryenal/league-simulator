package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"league-simulator/backend/db"
)

// MatchUpdateRequest represents the expected JSON body for updating a match score.
type MatchUpdateRequest struct {
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
}

// UpdateMatch handles PUT /match/{id}.
// It allows the result of a match to be updated manually.
func UpdateMatch(w http.ResponseWriter, r *http.Request) {
	// Ensure HTTP method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract match ID from the URL path
	idStr := strings.TrimPrefix(r.URL.Path, "/match/")
	matchID, err := strconv.Atoi(idStr)
	if err != nil || matchID <= 0 {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	// Parse JSON request body into a MatchUpdateRequest struct
	var update MatchUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Compute the result string based on scores
	var result string
	if update.HomeScore > update.AwayScore {
		result = "win"
	} else if update.HomeScore < update.AwayScore {
		result = "loss"
	} else {
		result = "draw"
	}

	// Prepare and execute the SQL update query
	stmt, err := db.DB.Prepare(`
		UPDATE matches 
		SET home_score = ?, away_score = ?, result = ?
		WHERE id = ?
	`)
	if err != nil {
		http.Error(w, "Failed to prepare update", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(update.HomeScore, update.AwayScore, result, matchID)
	if err != nil {
		http.Error(w, "Failed to execute update", http.StatusInternalServerError)
		return
	}

	// Confirm that the update affected exactly one row
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	// Return a simple success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Match %d updated: %d - %d (%s)", matchID, update.HomeScore, update.AwayScore, result)
}
