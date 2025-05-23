import React from 'react';
import '../styles/ChampionshipOdds.css';

/**
 * ChampionshipOdds component
 * Displays each team's predicted chance of winning the league.
 *
 * Props:
 * - odds: Array of objects with fields:
 *   - team (string)
 *   - chance (number, float) – percentage value between 0–100
 */
function ChampionshipOdds({ odds }) {
  if (!odds || odds.length === 0) {
    return (
      <div className="championship-odds">
        <p>No championship predictions available.</p>
      </div>
    );
  }

  return (
    <div className="championship-odds">
      <h2 className="championship-title">Championship Odds</h2>
      <table>
        <thead>
          <tr>
            <th>Team</th>
            <th>Chance (%)</th>
          </tr>
        </thead>
        <tbody>
          {odds.map((entry, index) => (
            <tr key={index}>
              <td>{entry.team}</td>
              <td>{entry.chance.toFixed(2)}%</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default ChampionshipOdds;
