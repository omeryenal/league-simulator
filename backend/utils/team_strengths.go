package utils

// TeamStrengths defines the base strength value for each team.
// These are used in simulation and prediction logic to determine match outcomes.
// Values are arbitrary but reflect real-world performance tiers.
var TeamStrengths = map[string]int{
	"Manchester City": 90, // Strongest team
	"Liverpool":       88,
	"Arsenal":         85,
	"Chelsea":         83,
}
