package components

import (
	"strings"
	"testing"
	"time"
)

func Test_makeDisplayName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"HighTowerLevel", args{"HighTowerLevel"}, "High Tower Level"},
		{"TowersBuilt", args{"TowersBuilt"}, "Towers Built"},
		{"Gold", args{"Gold"}, "Gold"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeDisplayName(tt.args.name); got != tt.want {
				t.Errorf("makeDisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGameStatsPreservesHighs(t *testing.T) {
	old := NewGameStats(nil)
	old.UpdateHighs(100, 4, 7)
	old.UpdateStat("TowersBuilt", 3)
	old.GameTime = 12 * time.Second

	got := NewGameStats(old)
	if got.GetStat("HighScore") != 100 {
		t.Errorf("HighScore = %v, want 100", got.GetStat("HighScore"))
	}
	if got.GetStat("HighCreepLevel") != 4 {
		t.Errorf("HighCreepLevel = %v, want 4", got.GetStat("HighCreepLevel"))
	}
	if got.GetStat("HighTowerLevel") != 7 {
		t.Errorf("HighTowerLevel = %v, want 7", got.GetStat("HighTowerLevel"))
	}
	if got.GetStat("TowersBuilt") != 0 {
		t.Errorf("TowersBuilt = %v, want 0 for a fresh run", got.GetStat("TowersBuilt"))
	}
	if got.GameTime != 0 {
		t.Errorf("GameTime = %v, want 0 for a fresh run", got.GameTime)
	}
}

func TestGameStatsUpdateAggregatesRunStats(t *testing.T) {
	total := NewGameStats(nil)
	total.UpdateHighs(100, 3, 6)
	total.UpdateStat("TowersBuilt", 2)
	total.GameTime = 10 * time.Second

	run := NewGameStats(nil)
	run.UpdateHighs(150, 2, 5)
	run.UpdateStat("TowersBuilt", 4)
	run.UpdateStat("CreepsKilled", 9)
	run.GameTime = 5 * time.Second

	total.Update(run)

	if got := total.GetStat("Games"); got != 1 {
		t.Errorf("Games = %v, want 1", got)
	}
	if got := total.GetStat("HighScore"); got != 150 {
		t.Errorf("HighScore = %v, want 150", got)
	}
	if got := total.GetStat("HighCreepLevel"); got != 3 {
		t.Errorf("HighCreepLevel = %v, want 3", got)
	}
	if got := total.GetStat("HighTowerLevel"); got != 6 {
		t.Errorf("HighTowerLevel = %v, want 6", got)
	}
	if got := total.GetStat("TowersBuilt"); got != 6 {
		t.Errorf("TowersBuilt = %v, want 6", got)
	}
	if got := total.GetStat("CreepsKilled"); got != 9 {
		t.Errorf("CreepsKilled = %v, want 9", got)
	}
	if total.GameTime != 15*time.Second {
		t.Errorf("GameTime = %v, want 15s", total.GameTime)
	}
}

func TestGameStatsResetPreservesHighsAndGames(t *testing.T) {
	stats := NewGameStats(nil)
	stats.UpdateHighs(100, 3, 6)
	stats.UpdateStat("Games", 2)
	stats.UpdateStat("TowersBuilt", 5)
	stats.GameTime = 10 * time.Second

	stats.Reset()

	if got := stats.GetStat("HighScore"); got != 100 {
		t.Errorf("HighScore = %v, want 100", got)
	}
	if got := stats.GetStat("HighCreepLevel"); got != 3 {
		t.Errorf("HighCreepLevel = %v, want 3", got)
	}
	if got := stats.GetStat("HighTowerLevel"); got != 6 {
		t.Errorf("HighTowerLevel = %v, want 6", got)
	}
	if got := stats.GetStat("Games"); got != 2 {
		t.Errorf("Games = %v, want 2", got)
	}
	if got := stats.GetStat("TowersBuilt"); got != 0 {
		t.Errorf("TowersBuilt = %v, want 0", got)
	}
	if stats.GameTime != 0 {
		t.Errorf("GameTime = %v, want 0", stats.GameTime)
	}
}

func TestGameStatsLinesFormatsStorageAndDisplayNames(t *testing.T) {
	stats := NewGameStats(nil)
	stats.UpdateStat("TowersBuilt", 2)
	stats.GameTime = 90 * time.Second

	storage := stats.StatsLines("=", false, false)
	if !strings.Contains(storage, "TowersBuilt=2\n") {
		t.Fatalf("storage stats missing TowersBuilt=2: %q", storage)
	}
	if !strings.Contains(storage, "GameTime=1m30s\n") {
		t.Fatalf("storage stats missing rounded GameTime: %q", storage)
	}

	display := stats.StatsLines(" ", false, true)
	if !strings.Contains(display, "Towers Built 2\n") {
		t.Fatalf("display stats missing display name: %q", display)
	}
}
func TestIterExcludeKeepsInputOrderAndSkipsExcludedValues(t *testing.T) {
	values := []string{"TowersBuilt", "Games", "HighScore", "CreepsKilled"}
	var got []string

	iterExclude(values, func(name string) {
		got = append(got, name)
	}, func(name string) bool {
		return name == "Games" || strings.HasPrefix(name, "High")
	})

	want := []string{"TowersBuilt", "CreepsKilled"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("iterExclude visited %v, want %v", got, want)
	}
}

func TestIterExcludePrefixSkipsAnyMatchingPrefix(t *testing.T) {
	values := []string{"BulletsExpired", "Games", "HighScore", "HighTowerLevel", "TowersBuilt"}
	var got []string

	iterExcludePrefix(values, func(name string) {
		got = append(got, name)
	}, "Game", "High")

	want := []string{"BulletsExpired", "TowersBuilt"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("iterExcludePrefix visited %v, want %v", got, want)
	}
}

func TestGameStatsIterExcludeVisitsSortedMapKeys(t *testing.T) {
	stats := NewGameStats(nil)
	stats.UpdateStat("TowersBuilt", 1)
	stats.UpdateStat("CreepsKilled", 2)
	stats.UpdateStat("Games", 3)
	stats.UpdateHighs(10, 2, 4)

	var got []string
	stats.iterExclude(func(name string) {
		got = append(got, name)
	}, func(name string) bool {
		return name == "Games" || strings.HasPrefix(name, "High")
	})

	want := []string{"CreepsKilled", "TowersBuilt"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("GameStats.iterExclude visited %v, want %v", got, want)
	}
}

func TestGameStatsIterExcludePrefixUsesSortedMapKeys(t *testing.T) {
	stats := NewGameStats(nil)
	stats.UpdateStat("TowersBuilt", 1)
	stats.UpdateStat("CreepsKilled", 2)
	stats.UpdateStat("Games", 3)
	stats.UpdateHighs(10, 2, 4)

	var got []string
	stats.iterExcludePrefix(func(name string) {
		got = append(got, name)
	}, "Game", "High")

	want := []string{"CreepsKilled", "TowersBuilt"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("GameStats.iterExcludePrefix visited %v, want %v", got, want)
	}
}
