import React, { useEffect, useState } from 'react';
import '../styles/MatchResults.css';

function MatchResults({ week }) {
  const [results, setResults] = useState([]);
  const [editingMatch, setEditingMatch] = useState(null);
  const [newScores, setNewScores] = useState({ home_score: 0, away_score: 0 });

  const API_URL = "http://localhost:8080"; // Manuel olarak sabitlendi

  useEffect(() => {
    if (!week) return;

    fetch(`${API_URL}/results/week/${week}`)
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) {
          setResults(data);
        } else {
          setResults([]);
        }
      })
      .catch(err => {
        console.error("Failed to fetch match results:", err);
        setResults([]);
      });
  }, [week]);

  const openModal = (match) => {
    setEditingMatch(match);
    setNewScores({ home_score: match.home_score, away_score: match.away_score });
  };

  const closeModal = () => {
    setEditingMatch(null);
  };

  const handleScoreChange = (e) => {
    const { name, value } = e.target;
    setNewScores(prev => ({ ...prev, [name]: parseInt(value) }));
  };

  const submitEdit = () => {
    fetch(`${API_URL}/match/${editingMatch.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newScores),
    })
      .then(res => {
        if (!res.ok) {
          throw new Error("Failed to update match");
        }
        setEditingMatch(null);
        window.location.reload();
      })
      .catch(err => {
        alert("Error updating match.");
        console.error(err);
      });
  };

  if (!results || results.length === 0) return <div>No matches found for week {week}</div>;

  return (
    <div className="match-results">
      <h2 className="match-title">Match Results - Week {week}</h2>
      <table>
        <thead>
          <tr>
            <th>Home Team</th>
            <th>Score</th>
            <th>Away Team</th>
            <th>Edit</th>
          </tr>
        </thead>
        <tbody>
          {results.map(match => (
            <tr key={match.id}>
              <td>{match.home_team}</td>
              <td>{match.home_score} - {match.away_score}</td>
              <td>{match.away_team}</td>
              <td>
                <button className="nav-button" onClick={() => openModal(match)}>Edit</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {editingMatch && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <h3>Edit Match</h3>
            <p>{editingMatch.home_team} vs {editingMatch.away_team}</p>
            <div>
              <input
                type="number"
                name="home_score"
                value={newScores.home_score}
                onChange={handleScoreChange}
                min="0"
              />
              {" - "}
              <input
                type="number"
                name="away_score"
                value={newScores.away_score}
                onChange={handleScoreChange}
                min="0"
              />
            </div>
            <br />
            <button className="nav-button" onClick={submitEdit}>Save</button>
            <button className="nav-button" onClick={closeModal} style={{ marginLeft: '10px' }}>Cancel</button>
          </div>
        </div>
      )}
    </div>
  );
}

export default MatchResults;
