import React, { useState } from 'react';

function EditMatchModal({ match, onClose, onSave }) {
  const [homeScore, setHomeScore] = useState(match.home_score);
  const [awayScore, setAwayScore] = useState(match.away_score);

  const API_URL = "http://localhost:8080"; // Manuel sabit URL

  const handleSubmit = async (e) => {
    e.preventDefault();

    const response = await fetch(`${API_URL}/match/${match.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ home_score: homeScore, away_score: awayScore }),
    });

    if (response.ok) {
      onSave();
      onClose();
    } else {
      alert("Failed to update match.");
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>Edit Match</h2>
        <form onSubmit={handleSubmit}>
          <label>
            {match.home_team}
            <input
              type="number"
              value={homeScore}
              onChange={(e) => setHomeScore(Number(e.target.value))}
              min="0"
            />
          </label>
          <label>
            {match.away_team}
            <input
              type="number"
              value={awayScore}
              onChange={(e) => setAwayScore(Number(e.target.value))}
              min="0"
            />
          </label>
          <button type="submit">Save</button>
          <button type="button" onClick={onClose} style={{ marginLeft: '10px' }}>Cancel</button>
        </form>
      </div>
    </div>
  );
}

export default EditMatchModal;
