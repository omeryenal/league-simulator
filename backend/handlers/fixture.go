package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"league-simulator/backend/db"
	"league-simulator/backend/models"
	"league-simulator/backend/utils"
)

// GetFixture handles the /fixture endpoint.
// It generates a fixture for N weeks and returns the schedule as JSON.
func GetFixture(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve all teams from the database
	rows, err := db.DB.Query("SELECT id, name FROM teams")
	if err != nil {
		http.Error(w, "Failed to fetch teams", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			http.Error(w, "Failed to scan team", http.StatusInternalServerError)
			return
		}
		teams = append(teams, t)
	}

	// We expect exactly 4 teams to generate the fixture
	if len(teams) != 4 {
		http.Error(w, "Fixture generation requires exactly 4 teams", http.StatusBadRequest)
		return
	}

	// Read the number of weeks from query parameter (?weeks=)
	weekParam := r.URL.Query().Get("weeks")
	weekCount := 12 // default
	if weekParam != "" {
		parsed, err := strconv.Atoi(weekParam)
		if err == nil && parsed > 0 {
			weekCount = parsed
		}
	}

	// Use a fixture generator service to generate the schedule
	generator := utils.NewSimpleFixtureService()
	fixture := generator.GenerateFixture(teams, weekCount)

	// Return the fixture as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fixture)
}
