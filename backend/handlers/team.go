package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"league-simulator/backend/db"
	"league-simulator/backend/models"
)

// CreateTeam handles POST /team.
// It adds a new team to the database using the name provided in the request body.
func CreateTeam(w http.ResponseWriter, r *http.Request) {
	// Ensure request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body to extract team name
	var team models.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil || team.Name == "" {
		http.Error(w, "Invalid team data", http.StatusBadRequest)
		return
	}

	// Insert the new team into the database
	stmt, err := db.DB.Prepare("INSERT INTO teams(name) VALUES(?)")
	if err != nil {
		http.Error(w, "Database error while preparing insert statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(team.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert team: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Team added successfully: %s", team.Name)
}

// GetTeams handles GET /teams.
// It returns all registered teams from the database as a JSON array.
func GetTeams(w http.ResponseWriter, r *http.Request) {
	// Ensure request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query all teams from the database
	rows, err := db.DB.Query("SELECT id, name FROM teams")
	if err != nil {
		http.Error(w, "Failed to fetch teams from the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teams []models.Team

	// Read and map each row to the Team struct
	for rows.Next() {
		var team models.Team
		err := rows.Scan(&team.ID, &team.Name)
		if err != nil {
			http.Error(w, "Failed to scan team row", http.StatusInternalServerError)
			return
		}
		teams = append(teams, team)
	}

	// Return the list of teams as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}
