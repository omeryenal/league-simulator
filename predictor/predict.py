import sys
import json
import random

# Constant boost for home team advantage
HOME_ADVANTAGE = 5

def predict_score(home_strength, away_strength):
    """
    Legacy function (unused): Predict goals using Gaussian sampling.
    Currently unused – kept for future extension.
    """
    diff = home_strength - away_strength
    base = 1.5 + diff * 0.03

    home_goals = max(0, int(random.gauss(base, 1)))
    away_goals = max(0, int(random.gauss(1.5, 1)))

    return home_goals, away_goals

def main():
    """
    Reads JSON match data from stdin, calculates predicted scores using a weighted power model,
    and outputs the result in JSON format.

    Each team's 'strength' is the primary input.
    Optional 'gd' (goal difference) adds a small impact to recent form.
    """
    input_data = json.load(sys.stdin)
    results = []

    for match in input_data:
        home = match['home_team']
        away = match['away_team']

        home_strength = home['strength']
        away_strength = away['strength']

        # Goal difference bonus (optional form boost)
        home_gd = home.get('gd', 0)
        away_gd = away.get('gd', 0)

        # Cap GD impact to avoid excessive swings
        gd_bonus = max(min(home_gd - away_gd, 5), -5)

        # Effective power: strength + GD bonus + home field advantage
        home_power = home_strength + gd_bonus + HOME_ADVANTAGE
        away_power = away_strength

        total_power = home_power + away_power

        # Normalize scores to a 3-goal scale
        if total_power == 0:
            home_score = away_score = 0
        else:
            home_score = round(3 * home_power / total_power)
            away_score = round(3 * away_power / total_power)

        # Avoid 0–0 outcomes: always give 1 goal to break ties
        if home_score == away_score:
            home_score += 1

        results.append({
            "home_team_id": home['id'],
            "away_team_id": away['id'],
            "home_score": home_score,
            "away_score": away_score
        })

    print(json.dumps(results))

if __name__ == "__main__":
    main()
