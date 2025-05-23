package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var DB *sql.DB

// InitDB opens a connection to the SQLite database and creates the necessary tables.
// Also inserts default teams and initial week 4 matches if they are missing.
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./league.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// SQL statements for creating teams and matches tables
	createTeamTable := `
	CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	`

	createMatchTable := `
	CREATE TABLE IF NOT EXISTS matches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		week INTEGER NOT NULL,
		home_team_id INTEGER,
		away_team_id INTEGER,
		home_score INTEGER,
		away_score INTEGER,
		result TEXT,
		FOREIGN KEY (home_team_id) REFERENCES teams(id),
		FOREIGN KEY (away_team_id) REFERENCES teams(id)
	);
	`

	// Execute table creation
	_, err = DB.Exec(createTeamTable)
	if err != nil {
		log.Fatal("Failed to create teams table:", err)
	}

	_, err = DB.Exec(createMatchTable)
	if err != nil {
		log.Fatal("Failed to create matches table:", err)
	}

	fmt.Println("Database connected and tables created successfully.")

	// Insert default teams and matches if necessary
	initTeams()
	initWeek4Matches()
}

// initTeams inserts the initial set of teams if the table is empty.
func initTeams() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM teams").Scan(&count)
	if err != nil {
		log.Println("Error checking teams count:", err)
		return
	}

	if count > 0 {
		return // Skip if teams already exist
	}

	_, err = DB.Exec(`
		INSERT INTO teams (name) VALUES
		('Manchester City'),
		('Liverpool'),
		('Arsenal'),
		('Chelsea')
	`)
	if err != nil {
		log.Println("Failed to insert teams:", err)
	} else {
		log.Println("Teams inserted successfully.")
	}
}

// initWeek4Matches adds results for week 4 if they aren't already present.
// These matches serve as a starting point for simulation.
func initWeek4Matches() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM matches WHERE week = 4").Scan(&count)
	if err != nil {
		log.Println("Error checking week 4 matches count:", err)
		return
	}

	if count > 0 {
		return // Matches already inserted
	}

	// Match data assumes teams have IDs 1 to 4 in the order they were inserted
	_, err = DB.Exec(`
		INSERT INTO matches (week, home_team_id, away_team_id, home_score, away_score, result) VALUES
		(4, 1, 2, 0, 0, 'draw'),
		(4, 3, 4, 1, 2, 'loss')
	`)
	if err != nil {
		log.Println("Failed to insert week 4 matches:", err)
	} else {
		log.Println("Week 4 matches inserted successfully.")
	}
}
