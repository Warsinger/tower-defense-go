//go:build !js

package components

import (
	"testing"
	"time"
)

func TestStatsStorageRoundTrip(t *testing.T) {
	t.Chdir(t.TempDir())
	stats := NewGameStats(nil)
	stats.UpdateHighs(1234, 7, 9)
	stats.UpdateStat("Games", 3)
	stats.GameTime = 42 * time.Second
	if err := stats.SaveStats(); err != nil {
		t.Fatal(err)
	}
	loaded := LoadStats()
	for _, name := range []string{"HighScore", "HighCreepLevel", "HighTowerLevel", "Games"} {
		if got, want := loaded.GetStat(name), stats.GetStat(name); got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}
	if loaded.GameTime != stats.GameTime {
		t.Errorf("game time = %v, want %v", loaded.GameTime, stats.GameTime)
	}
}
