package components

import (
	"testing"
	"tower-defense/config"
)

func TestPlayerData_GetCreepLevel(t *testing.T) {
	tests := []struct {
		name        string
		towerLevels int
		want        int
	}{
		{"zero tower levels starts at level one", 0, 1},
		{"four tower levels stays level one", 4, 1},
		{"five tower levels reaches level two", 5, 2},
		{"nineteen tower levels reaches level four", 19, 4},
		{"twenty tower levels reaches level five", 20, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerData{TowerLevels: tt.towerLevels}
			if got := player.GetCreepLevel(config.DefaultBalance()); got != tt.want {
				t.Errorf("GetCreepLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerData_GetMaxTowerLevel(t *testing.T) {
	tests := []struct {
		name        string
		towerLevels int
		want        int
	}{
		{"zero tower levels uses initial max", 0, 5},
		{"nineteen tower levels keeps initial max", 19, 5},
		{"twenty tower levels adds one max level", 20, 6},
		{"forty tower levels adds two max levels", 40, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := &PlayerData{TowerLevels: tt.towerLevels}
			if got := player.GetMaxTowerLevel(config.DefaultBalance()); got != tt.want {
				t.Errorf("GetMaxTowerLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
