package handlers

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"

	"league-simulator/backend/db"
	"league-simulator/backend/models"
	"league-simulator/backend/utils"
)

// Hardcoded strength values for each team used in probability calculations.
var defaultStrengths = map[string]int{
	"Manchester City": 90,
	"Liverpool":       85,
	"Arsenal":         80,
	"Chelsea":         75,
}

// MatchPrediction holds the calculated win/draw/lose probabilities and betting-style odds.
type MatchPrediction struct {
	HomeTeam   string  `json:"home_team"`
	AwayTeam   string  `json:"away_team"`
	HomeWinPct float64 `json:"home_win_pct"`
	DrawPct    float64 `json:"draw_pct"`
	AwayWinPct float64 `json:"away_win_pct"`
	HomeOdds   float64 `json:"home_odds"`
	DrawOdds   float64 `json:"draw_odds"`
	AwayOdds   float64 `json:"away_odds"`
}

// ChampionshipOdds represents the likelihood of a team becoming champion.
type ChampionshipOdds struct {
	TeamName string  `json:"team"`
	Chance   float64 `json:"chance"`
}

// PredictionResponse bundles both types of predictions into one response.
type PredictionResponse struct {
	Championship []ChampionshipOdds `json:"championship_odds"`
	NextWeek     []MatchPrediction  `json:"next_week_predictions"`
}

// GetPredictions handles GET /predictions.
// It calculates both championship odds and win/draw/lose odds for next week's matches.
func GetPredictions(w http.ResponseWriter, r *http.Request) {
	teams := getTeams()
	standings := getStandings()

	if len(teams) == 0 || len(standings) == 0 {
		http.Error(w, "Failed to retrieve teams or standings", http.StatusInternalServerError)
		return
	}

	// Championship odds calculation based on form and team strength
	var total float64
	weights := map[string]float64{}

	for _, s := range standings {
		strength, ok := defaultStrengths[s.TeamName]
		if !ok {
			strength = 70 // Default if team not found in map
		}
		score := float64(s.Wins*3+s.Draws) + float64(strength)*0.15
		weights[s.TeamName] = score
		total += score
	}

	var champOdds []ChampionshipOdds
	for team, weight := range weights {
		chance := math.Round((weight/total)*10000) / 100
		champOdds = append(champOdds, ChampionshipOdds{
			TeamName: team,
			Chance:   chance,
		})
	}

	// Determine the current week
	var currentWeek sql.NullInt64
	err := db.DB.QueryRow(`SELECT MAX(week) FROM matches`).Scan(&currentWeek)
	if err != nil || !currentWeek.Valid {
		http.Error(w, "Failed to determine current week", http.StatusInternalServerError)
		return
	}

	if currentWeek.Int64 >= 12 {
		// Season finished, no predictions to make
		json.NewEncoder(w).Encode(PredictionResponse{
			Championship: champOdds,
			NextWeek:     []MatchPrediction{},
		})
		return
	}

	// Get next week's fixture from the fixture generator
	weekIndex := int(currentWeek.Int64)
	fixture := utils.NewSimpleFixtureService().GenerateFixture(teams, 12)
	if weekIndex >= len(fixture) {
		http.Error(w, "Next week fixture not available", http.StatusBadRequest)
		return
	}

	// Match-by-match prediction with basic strength and past winner bonuses
	var weekPreds []MatchPrediction
	for _, match := range fixture[weekIndex] {
		home := match.HomeTeam
		away := match.AwayTeam

		homeStr := float64(defaultStrengths[home.Name])
		awayStr := float64(defaultStrengths[away.Name])

		// Past winner bonus: adds +4 strength
		winner := getPastWinner(home.ID, away.ID)
		if winner == home.ID {
			homeStr += 4
		} else if winner == away.ID {
			awayStr += 4
		}

		// Home advantage bonus
		homeStr += 5

		totalPower := homeStr + awayStr
		homePct := 100 * 0.75 * (homeStr / totalPower)
		awayPct := 100 * 0.75 * (awayStr / totalPower)
		drawPct := 25.0 // fixed draw chance

		// Odds = inverse probability
		homeOdds := math.Round((100/homePct)*100) / 100
		drawOdds := math.Round((100/drawPct)*100) / 100
		awayOdds := math.Round((100/awayPct)*100) / 100

		weekPreds = append(weekPreds, MatchPrediction{
			HomeTeam:   home.Name,
			AwayTeam:   away.Name,
			HomeWinPct: math.Round(homePct*100) / 100,
			DrawPct:    math.Round(drawPct*100) / 100,
			AwayWinPct: math.Round(awayPct*100) / 100,
			HomeOdds:   homeOdds,
			DrawOdds:   drawOdds,
			AwayOdds:   awayOdds,
		})
	}

	// Return final combined prediction response
	response := PredictionResponse{
		Championship: champOdds,
		NextWeek:     weekPreds,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getTeams queries the DB and returns a list of all teams.
func getTeams() []models.Team {
	rows, err := db.DB.Query("SELECT id, name FROM teams")
	if err != nil {
		return []models.Team{}
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			continue
		}
		teams = append(teams, t)
	}
	return teams
}

// getStandings calculates the current league standings from match results.
func getStandings() []models.Standing {
	rows, err := db.DB.Query(`
		SELECT 
			t.id, t.name,
			COUNT(m.id) as played,
			SUM(CASE WHEN (t.id = m.home_team_id AND m.result = 'win') OR (t.id = m.away_team_id AND m.result = 'loss') THEN 1 ELSE 0 END) AS wins,
			SUM(CASE WHEN m.result = 'draw' AND (t.id = m.home_team_id OR t.id = m.away_team_id) THEN 1 ELSE 0 END) AS draws,
			SUM(CASE WHEN (t.id = m.home_team_id AND m.result = 'loss') OR (t.id = m.away_team_id AND m.result = 'win') THEN 1 ELSE 0 END) AS losses,
			SUM(CASE WHEN (t.id = m.home_team_id AND m.result = 'win') OR (t.id = m.away_team_id AND m.result = 'loss') THEN 3 ELSE 0 END) +
			SUM(CASE WHEN m.result = 'draw' AND (t.id = m.home_team_id OR t.id = m.away_team_id) THEN 1 ELSE 0 END) AS points
		FROM teams t
		LEFT JOIN matches m ON t.id = m.home_team_id OR t.id = m.away_team_id
		GROUP BY t.id
		ORDER BY points DESC, wins DESC
	`)
	if err != nil {
		return []models.Standing{}
	}
	defer rows.Close()

	var standings []models.Standing
	for rows.Next() {
		var s models.Standing
		if err := rows.Scan(&s.TeamID, &s.TeamName, &s.Played, &s.Wins, &s.Draws, &s.Losses, &s.Points); err != nil {
			continue
		}
		standings = append(standings, s)
	}
	return standings
}

// getPastWinner checks the last match result between two teams and returns the winner's team ID.
// If draw or no history, returns 0.
func getPastWinner(id1, id2 int) int {
	row := db.DB.QueryRow(`
		SELECT home_team_id, away_team_id, result 
		FROM matches 
		WHERE (home_team_id = ? AND away_team_id = ?) 
		   OR (home_team_id = ? AND away_team_id = ?)
		ORDER BY id DESC LIMIT 1
	`, id1, id2, id2, id1)

	var hID, aID int
	var res string
	err := row.Scan(&hID, &aID, &res)
	if err != nil {
		return 0
	}

	if res == "draw" {
		return 0
	}

	if (hID == id1 && res == "win") || (aID == id1 && res == "loss") {
		return id1
	}
	return id2
}
