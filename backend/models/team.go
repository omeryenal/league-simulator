package models

// Team represents a football team in the league.
type Team struct {
	ID   int    `json:"id"`   // Unique identifier for the team
	Name string `json:"name"` // Display name of the team
}
