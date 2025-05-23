import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';

import './App.css';
import HomePage from './pages/HomePage';
import PredictionsPage from './pages/PredictionsPage';

/**
 * Navigation component
 * Renders navigation buttons depending on current path and league week.
 * - Home button: visible if not already on homepage
 * - Predictions button: visible if week <= 12 (league not finished)
 */
function Navigation() {
  const location = useLocation();
  const [week, setWeek] = useState(1);

  // Fetch current league week to decide whether predictions should be available
  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/week/current`)
      .then(res => res.json())
      .then(({ week }) => setWeek(week))
      .catch(console.error);
  }, []);
  

  return (
    <nav className="navbar">
      {/* Show Home button if not on the homepage */}
      {location.pathname !== '/' && (
        <Link to="/">
          <button className="nav-button">Home</button>
        </Link>
      )}
      {/* Show Predictions button only if league is ongoing */}
      {location.pathname !== '/predictions' && week <= 12 && (
        <Link to="/predictions">
          <button className="nav-button">Predictions</button>
        </Link>
      )}
    </nav>
  );
}

/**
 * App component
 * Wraps the entire application with router and renders page-level routes.
 */
function App() {
  return (
    <Router>
      <div className="App">
        <header>
          <h1 className="main-title">⚽ League Simulator</h1>
          <Navigation />
        </header>

        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/predictions" element={<PredictionsPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
