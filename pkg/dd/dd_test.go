package dd

import (
	"testing"
)

func testGame(t *testing.T) {
	// Create a new game
	g := NewGame()
	wantPlayers := map[string]bool{"alice": true, "bob": true, "charlie": true, "david": true}
	teams := map[string][]string{"team1": {"alice", "charlie"}, "team2": {"david", "jose"}}
	for team, players := range teams {
		t := NewTeam(team)
		for _, player := range players {
			p := NewPlayer(player)
			t.AddPlayer(p)
		}
		g.AddTeam(t)
	}
	if len(g.teams) != 2 {
		t.Errorf("Expected 2 teams, got %d", len(g.teams))
	}
	if len(g.players) != 4 {
		t.Errorf("Expected 4 players, got %d", len(g.players))
	}
	for _, p := range g.players {
		if !wantPlayers[p.name] {
			t.Errorf("Unexpected player %s", p.name)
		}
	}
}
