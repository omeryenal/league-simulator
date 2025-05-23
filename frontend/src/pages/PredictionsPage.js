import React, { useEffect, useState } from 'react';
import ChampionshipOdds from '../components/ChampionshipOdds';
import BettingOdds from '../components/BettingOdds';

function PredictionsPage() {
  const [champOdds, setChampOdds] = useState([]);
  const [bettingOdds, setBettingOdds] = useState([]);

  const API_URL = process.env.REACT_APP_API_URL;

  useEffect(() => {
    fetch(`${API_URL}/predictions`)
      .then((res) => res.json())
      .then((data) => {
        setChampOdds(data.championship_odds);
        setBettingOdds(data.next_week_predictions);
      })
      .catch((error) => {
        console.error("Error fetching predictions:", error);
      });
  }, [API_URL]);

  return (
    <>
      <ChampionshipOdds odds={champOdds} />
      <BettingOdds matches={bettingOdds} />
    </>
  );
}

export default PredictionsPage;
