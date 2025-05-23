package utils

import (
	"league-simulator/backend/models"
)

// MatchPair represents a scheduled match between two teams.
type MatchPair struct {
	HomeTeam models.Team `json:"home_team"`
	AwayTeam models.Team `json:"away_team"`
}

// FixtureGenerator is an interface for any service that can generate a fixture.
// You can implement different algorithms (e.g. round-robin, random shuffle) behind this.
type FixtureGenerator interface {
	GenerateFixture(teams []models.Team, weekCount int) [][]MatchPair
}

// SimpleFixtureService is a round-robin based fixture generator.
// It creates a balanced schedule where each team plays every other team.
type SimpleFixtureService struct{}

// NewSimpleFixtureService returns a new instance that satisfies FixtureGenerator.
func NewSimpleFixtureService() FixtureGenerator {
	return &SimpleFixtureService{}
}

// GenerateFixture creates a round-robin style fixture for a given list of teams.
// If weekCount > total unique rounds, the fixture is repeated.
func (s *SimpleFixtureService) GenerateFixture(teams []models.Team, weekCount int) [][]MatchPair {
	n := len(teams)

	// If number of teams is odd, add a dummy team (BYE) to make it even
	if n%2 != 0 {
		teams = append(teams, models.Team{ID: 0, Name: "BYE"})
		n++
	}

	totalRounds := n - 1
	matchesPerRound := n / 2
	weeks := make([][]MatchPair, totalRounds)

	// Generate matches using round-robin pairing
	for round := 0; round < totalRounds; round++ {
		var week []MatchPair

		for match := 0; match < matchesPerRound; match++ {
			homeIdx := (round + match) % (n - 1)
			awayIdx := (n - 1 - match + round) % (n - 1)

			// The last team is fixed in rotation
			if match == 0 {
				awayIdx = n - 1
			}

			homeTeam := teams[homeIdx]
			awayTeam := teams[awayIdx]

			// Skip BYE matches
			if homeTeam.ID == 0 || awayTeam.ID == 0 {
				continue
			}

			// Alternate home/away for fairness
			if round%2 == 0 {
				week = append(week, MatchPair{HomeTeam: homeTeam, AwayTeam: awayTeam})
			} else {
				week = append(week, MatchPair{HomeTeam: awayTeam, AwayTeam: homeTeam})
			}
		}

		weeks[round] = week
	}

	// Repeat the fixture if more weeks are needed
	for len(weeks) < weekCount {
		weeks = append(weeks, weeks[len(weeks)%totalRounds])
	}

	return weeks[:weekCount]
}
