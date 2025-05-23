package handlers

import (
	"encoding/json"
	"net/http"

	"league-simulator/backend/db"
	"league-simulator/backend/models"
)

// GetStandings handles GET /standings.
// It returns the league table with points, goal difference, and other metrics for each team.
func GetStandings(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// SQL query to compute team standings
	// Includes matches played, wins, draws, losses, goal difference, and points
	query := `
	SELECT 
		t.id AS team_id,
		t.name AS team_name,
		COUNT(m.id) AS played,
		SUM(CASE WHEN (t.id = m.home_team_id AND m.result = 'win') OR (t.id = m.away_team_id AND m.result = 'loss') THEN 1 ELSE 0 END) AS wins,
		SUM(CASE WHEN m.result = 'draw' AND (t.id = m.home_team_id OR t.id = m.away_team_id) THEN 1 ELSE 0 END) AS draws,
		SUM(CASE WHEN (t.id = m.home_team_id AND m.result = 'loss') OR (t.id = m.away_team_id AND m.result = 'win') THEN 1 ELSE 0 END) AS losses,
		SUM(
			CASE 
				WHEN t.id = m.home_team_id THEN m.home_score - m.away_score
				WHEN t.id = m.away_team_id THEN m.away_score - m.home_score
				ELSE 0 
			END
		) AS goal_difference,
		SUM(CASE WHEN (t.id = m.home_team_id AND m.result = 'win') OR (t.id = m.away_team_id AND m.result = 'loss') THEN 3 ELSE 0 END) +
		SUM(CASE WHEN m.result = 'draw' AND (t.id = m.home_team_id OR t.id = m.away_team_id) THEN 1 ELSE 0 END) AS points
	FROM teams t
	LEFT JOIN matches m ON t.id = m.home_team_id OR t.id = m.away_team_id
	GROUP BY t.id
	ORDER BY points DESC, goal_difference DESC, wins DESC
	`

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to calculate standings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Prepare the response slice
	var standings []models.Standing
	for rows.Next() {
		var s models.Standing
		err := rows.Scan(
			&s.TeamID,
			&s.TeamName,
			&s.Played,
			&s.Wins,
			&s.Draws,
			&s.Losses,
			&s.GoalDifference,
			&s.Points,
		)
		if err != nil {
			http.Error(w, "Failed to scan standings row", http.StatusInternalServerError)
			return
		}
		standings = append(standings, s)
	}

	// Return the standings as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(standings)
}
