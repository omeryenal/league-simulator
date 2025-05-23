import React from 'react';
import '../styles/LeagueTable.css'; // Table styling

/**
 * LeagueTable component
 * Displays current league standings in a table format.
 *
 * Props:
 * - standings: Array of objects with fields:
 *   - team_name, points, played, wins, draws, losses, goal_difference
 */
function LeagueTable({ standings }) {
  if (!standings || standings.length === 0) {
    return (
      <div className="league-table">
        <p>No standings data available.</p>
      </div>
    );
  }

  return (
    <div className="league-table">
      <h2 className="league-title">League Standings</h2>
      <table className="league-table">
        <thead>
          <tr>
            <th>Team</th>
            <th>PTS</th> {/* Points */}
            <th>P</th>   {/* Played */}
            <th>W</th>   {/* Wins */}
            <th>D</th>   {/* Draws */}
            <th>L</th>   {/* Losses */}
            <th>GD</th>  {/* Goal Difference */}
          </tr>
        </thead>
        <tbody>
          {standings.map((team, index) => (
            <tr key={index}>
              <td>{team.team_name}</td>
              <td>{team.points}</td>
              <td>{team.played}</td>
              <td>{team.wins}</td>
              <td>{team.draws}</td>
              <td>{team.losses}</td>
              <td>{team.goal_difference}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default LeagueTable;
