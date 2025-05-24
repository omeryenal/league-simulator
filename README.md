# ⚽ League Simulator

A full-stack football league simulation platform that lets you simulate matches, view weekly results, edit scores, track standings, and get betting-style predictions — all in real-time.

> Built with Go (backend), Python (prediction logic), and React (frontend).

---

## 📌 Features

- 🎲 Simulate weekly matches or the entire season
- ✏️ Manually edit match results at any time
- 📈 Real-time standings with wins, draws, losses, goal difference, and points
- 🧠 Python-powered match prediction engine
- 🔢 Betting-style win/draw/loss odds based on strength, form, and history

---

## 🧠 Architecture

```txt
User Browser (React Frontend)
         ↓
Localhost API (Port 3000)
         ↓
Go Backend (Port 8080)
       ↙        ↘
 SQLite DB     Python Prediction Engine


🛠 Tech Stack

Layer	Tech
Frontend	React, HTML/CSS, Tailwind
Backend	Go (net/http), SQLite
ML Logic	Python (match prediction script)
Database	SQLite (local storage)
⚙️ Local Setup

1. Clone the repo
git clone https://github.com/omeryenal/league-simulator.git
cd league-simulator
2. Run the backend
cd backend
go run main.go
3. Run the frontend
cd frontend
npm install
npm start
Make sure your .env file inside frontend/ contains:
REACT_APP_API_URL=http://localhost:8080

🧪 Sample Use Cases

✏️ Edit Match Results
Go to the Match Results page.
Click the Edit button next to any match.
Update the home_score or away_score fields.
Standings update instantly after saving.

⚽ Simulate Weekly Matches
Click "Simulate Week" to play the current week.
Each match is predicted using a Python script based on:
Team strength
Home advantage
Current form
Results and standings update automatically.

🏁 Play All Weeks
Click "Play All Weeks" to simulate the full season in one click.
All matches are played automatically.
Final standings and champion probabilities are shown.

🧮 Prediction & Odds Engine

Each match prediction includes:

Realistic win/draw/loss percentages
Bookmaker-style odds (calculated from implied probabilities)
Factoring in:
Past performance
Team strength
Head-to-head results
Example output:

Team A (Home) vs Team B
→  Team A Win: 48% (odds: 2.08)
→  Draw: 28% (odds: 3.57)
→  Team B Win: 24% (odds: 4.17)
📁 Folder Structure

league-simulator/
│
├── backend/              # Go server + database + API routes
│   ├── db/               # SQLite setup and queries
│   ├── handlers/         # HTTP handlers
│   └── main.go
│
├── frontend/             # React frontend
│   ├── src/components/   # League table, match list, edit modal
│   └── .env              # API URL
│
├── predict/              # Python match prediction logic
│   └── predict.py
🚀 Future Improvements

📊 View match history and stats
🔐 Add user authentication and private leagues
📅 Custom fixture generator
🧠 Advanced prediction models (e.g., XGBoost, LSTM)
🌍 Cloud deployment (Railway, Vercel, or Render)
🙋‍♂️ Author

Made by Ömer Yenal