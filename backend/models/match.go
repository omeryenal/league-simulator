package models

// Match represents a single match between two teams in a given week.
type Match struct {
	ID         int    `json:"id"`           // Unique ID of the match
	Week       int    `json:"week"`         // Week number when the match was played
	HomeTeamID int    `json:"home_team_id"` // ID of the home team
	AwayTeamID int    `json:"away_team_id"` // ID of the away team
	HomeScore  int    `json:"home_score"`   // Goals scored by home team
	AwayScore  int    `json:"away_score"`   // Goals scored by away team
	Result     string `json:"result"`       // Outcome from home team's perspective: "win", "loss", or "draw"
}
