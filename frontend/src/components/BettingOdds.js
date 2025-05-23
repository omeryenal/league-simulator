import React from 'react';
import '../styles/BettingOdds.css';

/**
 * BettingOdds component
 * Displays predicted win/draw/lose probabilities and betting odds for upcoming matches.
 *
 * Props:
 * - matches: Array of objects with the following fields:
 *   - home_team, away_team (string)
 *   - home_win_pct, draw_pct, away_win_pct (number, %)
 *   - home_odds, draw_odds, away_odds (float)
 */
function BettingOdds({ matches }) {
  if (!matches || matches.length === 0) {
    return (
      <div className="betting-odds">
        <p>No betting data available for the next week.</p>
      </div>
    );
  }

  return (
    <div className="betting-odds">
      <h2 className="betting-title">Betting Odds – Next Week</h2>
      <table>
        <thead>
          <tr>
            <th>Match</th>
            <th>Home Win %</th>
            <th>Draw %</th>
            <th>Away Win %</th>
            <th>Home Odds</th>
            <th>Draw Odds</th>
            <th>Away Odds</th>
          </tr>
        </thead>
        <tbody>
          {matches.map((match, index) => (
            <tr key={index}>
              <td>{match.home_team} vs {match.away_team}</td>
              <td>{match.home_win_pct}%</td>
              <td>{match.draw_pct}%</td>
              <td>{match.away_win_pct}%</td>
              <td>{match.home_odds.toFixed(2)}</td>
              <td>{match.draw_odds.toFixed(2)}</td>
              <td>{match.away_odds.toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default BettingOdds;
