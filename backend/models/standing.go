package models

// Standing represents the league table status of a team.
type Standing struct {
	TeamID         int    `json:"team_id"`         // Unique ID of the team
	TeamName       string `json:"team_name"`       // Name of the team
	Played         int    `json:"played"`          // Total number of matches played
	Wins           int    `json:"wins"`            // Number of wins
	Draws          int    `json:"draws"`           // Number of draws
	Losses         int    `json:"losses"`          // Number of losses
	GoalDifference int    `json:"goal_difference"` // Total goal difference (goals scored - goals conceded)
	Points         int    `json:"points"`          // Total points (win = 3 pts, draw = 1 pt, loss = 0 pts)
}
