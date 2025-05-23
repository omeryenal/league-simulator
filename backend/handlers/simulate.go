package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"league-simulator/backend/db"
	"league-simulator/backend/models"
	"league-simulator/backend/utils"
)

// TeamStrengths defines the base strength rating for each team.
// These values are used in the Python prediction script.
var TeamStrengths = map[string]int{
	"Manchester City": 85,
	"Liverpool":       83,
	"Arsenal":         78,
	"Chelsea":         75,
}

// simulateWeekAndInsert simulates the results of a given week using a Python script.
// It clears old matches for that week and stores the new simulated results in the DB.
func simulateWeekAndInsert(week int, matches []models.Match, teams []models.Team) ([]map[string]interface{}, error) {
	// Remove existing matches for this week to avoid duplicates
	_, err := db.DB.Exec("DELETE FROM matches WHERE week = ?", week)
	if err != nil {
		return nil, fmt.Errorf("Failed to clear old matches for week %d: %v", week, err)
	}

	// Map team IDs to their data for easy lookup
	teamMap := make(map[int]models.Team)
	for _, t := range teams {
		teamMap[t.ID] = t
	}

	// Build input payload for the Python script
	var input []map[string]interface{}
	for _, match := range matches {
		home := teamMap[match.HomeTeamID]
		away := teamMap[match.AwayTeamID]

		input = append(input, map[string]interface{}{
			"home_team": map[string]interface{}{
				"id":       home.ID,
				"name":     home.Name,
				"strength": TeamStrengths[home.Name],
			},
			"away_team": map[string]interface{}{
				"id":       away.ID,
				"name":     away.Name,
				"strength": TeamStrengths[away.Name],
			},
		})
	}

	// Convert input to JSON and run the Python prediction script
	jsonInput, _ := json.Marshal(input)
	cmd := exec.Command("python3", "../predictor/predict.py")
	cmd.Stdin = bytes.NewReader(jsonInput)
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Python error: %v", err)
	}

	// Decode output from the Python script
	var results []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &results); err != nil {
		return nil, fmt.Errorf("Invalid JSON output: %v", err)
	}

	// Prepare SQL insert statement
	stmt, err := db.DB.Prepare(`
		INSERT INTO matches (week, home_team_id, away_team_id, home_score, away_score, result)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("DB prepare error: %v", err)
	}
	defer stmt.Close()

	// Save each predicted match into the database
	for _, match := range results {
		home := int(match["home_score"].(float64))
		away := int(match["away_score"].(float64))

		result := "draw"
		if home > away {
			result = "win"
		} else if home < away {
			result = "loss"
		}

		_, err := stmt.Exec(
			week,
			int(match["home_team_id"].(float64)),
			int(match["away_team_id"].(float64)),
			home,
			away,
			result,
		)
		if err != nil {
			return nil, fmt.Errorf("DB insert error: %v", err)
		}
	}

	return results, nil
}

// SimulateWeek handles GET /simulate/week?n=5
// It simulates only the selected week and stores the result.
func SimulateWeek(w http.ResponseWriter, r *http.Request) {
	setupCORS(w, r)

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get week number from query string
	weekParam := r.URL.Query().Get("n")
	weekIndex, err := strconv.Atoi(weekParam)
	if err != nil || weekIndex < 1 {
		http.Error(w, "Invalid week index", http.StatusBadRequest)
		return
	}

	teams, err := fetchTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate full fixture and extract that specific week
	fixture := utils.NewSimpleFixtureService().GenerateFixture(teams, 14)
	if weekIndex > len(fixture) {
		http.Error(w, fmt.Sprintf("Week %d does not exist", weekIndex), http.StatusBadRequest)
		return
	}

	var weekMatches []models.Match
	for _, mp := range fixture[weekIndex-1] {
		weekMatches = append(weekMatches, models.Match{
			HomeTeamID: mp.HomeTeam.ID,
			AwayTeamID: mp.AwayTeam.ID,
		})
	}

	results, err := simulateWeekAndInsert(weekIndex, weekMatches, teams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

// SimulateNextWeek handles POST /simulate/next
// It simulates the next unplayed week based on the current progress.
func SimulateNextWeek(w http.ResponseWriter, r *http.Request) {
	setupCORS(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	const MaxWeek = 12

	row := db.DB.QueryRow(`SELECT MAX(week) FROM matches`)
	var lastPlayed sql.NullInt64
	err := row.Scan(&lastPlayed)
	if err != nil {
		http.Error(w, "Failed to get last played week", http.StatusInternalServerError)
		return
	}

	nextWeek := 1
	if lastPlayed.Valid {
		nextWeek = int(lastPlayed.Int64) + 1
	}

	if nextWeek > MaxWeek {
		http.Error(w, fmt.Sprintf("Week %d exceeds max week limit", nextWeek), http.StatusBadRequest)
		return
	}

	teams, err := fetchTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fixture := utils.NewSimpleFixtureService().GenerateFixture(teams, MaxWeek)
	if nextWeek > len(fixture) {
		http.Error(w, fmt.Sprintf("Week %d does not exist", nextWeek), http.StatusBadRequest)
		return
	}

	var weekMatches []models.Match
	for _, mp := range fixture[nextWeek-1] {
		weekMatches = append(weekMatches, models.Match{
			HomeTeamID: mp.HomeTeam.ID,
			AwayTeamID: mp.AwayTeam.ID,
		})
	}

	results, err := simulateWeekAndInsert(nextWeek, weekMatches, teams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

// SimulateAll handles POST /simulate/all
// It simulates the entire season from week 4 to weekCount.
func SimulateAll(w http.ResponseWriter, r *http.Request) {
	setupCORS(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	const weekCount = 12

	teams, err := fetchTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fixture := utils.NewSimpleFixtureService().GenerateFixture(teams, weekCount)
	allResults := make([][]map[string]interface{}, 0)

	for i, week := range fixture {
		weekNumber := i + 1

		// Clear old results if any exist for this week
		_, err := db.DB.Exec("DELETE FROM matches WHERE week = ?", weekNumber)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to clear old matches for week %d: %v", weekNumber, err), http.StatusInternalServerError)
			return
		}

		var weekMatches []models.Match
		for _, mp := range week {
			weekMatches = append(weekMatches, models.Match{
				HomeTeamID: mp.HomeTeam.ID,
				AwayTeamID: mp.AwayTeam.ID,
			})
		}

		results, err := simulateWeekAndInsert(weekNumber, weekMatches, teams)
		if err != nil {
			http.Error(w, fmt.Sprintf("Simulation failed on week %d: %v", weekNumber, err), http.StatusInternalServerError)
			return
		}

		allResults = append(allResults, results)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allResults)
}

// fetchTeams returns all teams from the database.
func fetchTeams() ([]models.Team, error) {
	rows, err := db.DB.Query("SELECT id, name FROM teams")
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch teams: %v", err)
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, fmt.Errorf("Failed to scan team: %v", err)
		}
		teams = append(teams, t)
	}
	return teams, nil
}

// setupCORS sets CORS headers for cross-origin requests.
func setupCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}
