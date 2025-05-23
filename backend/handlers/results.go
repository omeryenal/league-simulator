package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"league-simulator/backend/db"
)

// GetWeekResults handles GET /results/week/{n}.
// It returns all match results for a given week in JSON format.
func GetWeekResults(w http.ResponseWriter, r *http.Request) {
	// Enforce GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract week number from URL path
	weekStr := strings.TrimPrefix(r.URL.Path, "/results/week/")
	weekNum, err := strconv.Atoi(weekStr)
	if err != nil || weekNum <= 0 {
		http.Error(w, "Invalid week number", http.StatusBadRequest)
		return
	}

	// Query the database for matches played in the specified week
	rows, err := db.DB.Query(`
		SELECT m.id, m.week, t1.name, t2.name, m.home_score, m.away_score, m.result
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		WHERE m.week = ?
		ORDER BY m.id ASC
	`, weekNum)
	if err != nil {
		http.Error(w, "Failed to query matches", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// MatchResult defines how each match will appear in the response JSON
	type MatchResult struct {
		ID        int    `json:"id"`
		Week      int    `json:"week"`
		HomeTeam  string `json:"home_team"`
		AwayTeam  string `json:"away_team"`
		HomeScore int    `json:"home_score"`
		AwayScore int    `json:"away_score"`
		Result    string `json:"result"`
	}

	// Read all rows into a results slice
	var results []MatchResult
	for rows.Next() {
		var res MatchResult
		err := rows.Scan(&res.ID, &res.Week, &res.HomeTeam, &res.AwayTeam, &res.HomeScore, &res.AwayScore, &res.Result)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		results = append(results, res)
	}

	// Return results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
