package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"league-simulator/backend/db"
	"league-simulator/backend/models"
)

// CreateMatch handles POST /match.
// It inserts a new match into the database, including scores and result.
func CreateMatch(w http.ResponseWriter, r *http.Request) {
	// Ensure request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body into a Match struct
	var match models.Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		http.Error(w, "Invalid match data", http.StatusBadRequest)
		return
	}

	// Compute match result based on score
	if match.HomeScore > match.AwayScore {
		match.Result = "win"
	} else if match.HomeScore < match.AwayScore {
		match.Result = "loss"
	} else {
		match.Result = "draw"
	}

	// Prepare the SQL insert statement
	stmt, err := db.DB.Prepare(`
		INSERT INTO matches (week, home_team_id, away_team_id, home_score, away_score, result)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		http.Error(w, "Database prepare error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the statement with match data
	_, err = stmt.Exec(
		match.Week,
		match.HomeTeamID,
		match.AwayTeamID,
		match.HomeScore,
		match.AwayScore,
		match.Result,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert match: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with confirmation
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Match added successfully (Week %d)", match.Week)
}
