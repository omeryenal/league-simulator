import React, { useEffect, useState } from 'react';
import LeagueTable from '../components/LeagueTable';
import MatchResults from '../components/MatchResults';

function HomePage() {
  const [standings, setStandings] = useState([]);
  const [currentWeek, setCurrentWeek] = useState(4);
  const [playAllClicked, setPlayAllClicked] = useState(false);
  const MAX_WEEK = 12;

  const API_URL = process.env.REACT_APP_API_URL;

  // Fetch current week
  useEffect(() => {
    fetch(`${API_URL}/week/current`)
      .then(res => res.json())
      .then(({ week }) => setCurrentWeek(week))
      .catch(console.error);
  }, [API_URL]);

  // Fetch current standings
  useEffect(() => {
    fetch(`${API_URL}/standings`)
      .then(res => res.json())
      .then(setStandings)
      .catch(console.error);
  }, [API_URL]);

  // Simulate all weeks
  const handlePlayAll = async () => {
    await fetch(`${API_URL}/simulate/all`, { method: "POST" });
    setPlayAllClicked(true);
    window.location.reload();
  };

  // Simulate next week
  const handleNextWeek = async () => {
    await fetch(`${API_URL}/simulate/next`, { method: "POST" });
    window.location.reload();
  };

  // Reset season to week 4
  const handleResetSeason = async () => {
    await fetch(`${API_URL}/reset`, { method: "POST" });
    setPlayAllClicked(false);
    setCurrentWeek(4);
    window.location.reload();
  };

  const champion = standings.length > 0 ? standings[0].team_name : null;

  return (
    <>
      <LeagueTable standings={standings} />
      <MatchResults week={Math.min(currentWeek - 1, 12)} />

      {currentWeek > MAX_WEEK && champion && (
        <div className="champion-top-banner">
          🏆 Champion: <strong>{champion}</strong>
        </div>
      )}

      {currentWeek <= MAX_WEEK && !playAllClicked && (
        <button onClick={handlePlayAll} className="bottom-left-btn">Play All</button>
      )}

      {currentWeek <= MAX_WEEK && (
        <button onClick={handleNextWeek} className="bottom-right-btn">Next Week</button>
      )}

      <button
        onClick={handleResetSeason}
        style={{
          position: 'fixed',
          top: 20,
          right: 20,
          padding: '10px 20px',
          backgroundColor: '#ff4d4d',
          color: 'white',
          border: 'none',
          borderRadius: 8,
          fontWeight: 'bold',
          cursor: 'pointer',
          boxShadow: '0 2px 6px rgba(0,0,0,0.15)',
          fontFamily: 'Poppins, sans-serif',
          fontSize: 15,
          zIndex: 999
        }}
      >
        Reset Season
      </button>
    </>
  );
}

export default HomePage;
