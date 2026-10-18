package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBalanceUsesEmbeddedDefaultForEmptyPath(t *testing.T) {
	balance, err := LoadBalance("")
	if err != nil {
		t.Fatalf("LoadBalance(\"\") error = %v", err)
	}

	if balance.Player.StartingMoney != DefaultBalance().Player.StartingMoney {
		t.Errorf("StartingMoney = %v, want embedded default %v", balance.Player.StartingMoney, DefaultBalance().Player.StartingMoney)
	}
	if balance.Tower.Costs["Ranged"] != DefaultBalance().Tower.Costs["Ranged"] {
		t.Errorf("Ranged tower cost = %v, want embedded default %v", balance.Tower.Costs["Ranged"], DefaultBalance().Tower.Costs["Ranged"])
	}
}

func TestLoadBalanceReadsExternalJSON(t *testing.T) {
	want := *DefaultBalance()
	want.Player.StartingMoney = 1234
	want.Tower.DefaultType = "Laser"
	want.Tower.Costs = map[string]int{"Laser": 77}
	want.Multiplayer.SuperCreepCooldown = 42

	bytes, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	path := filepath.Join(t.TempDir(), "balance.json")
	if err := os.WriteFile(path, bytes, 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	got, err := LoadBalance(path)
	if err != nil {
		t.Fatalf("LoadBalance(%q) error = %v", path, err)
	}

	if got.Player.StartingMoney != 1234 {
		t.Errorf("StartingMoney = %v, want 1234", got.Player.StartingMoney)
	}
	if got.Tower.DefaultType != "Laser" {
		t.Errorf("DefaultType = %q, want Laser", got.Tower.DefaultType)
	}
	if got.Tower.Costs["Laser"] != 77 {
		t.Errorf("Laser tower cost = %v, want 77", got.Tower.Costs["Laser"])
	}
	if got.Multiplayer.SuperCreepCooldown != 42 {
		t.Errorf("SuperCreepCooldown = %v, want 42", got.Multiplayer.SuperCreepCooldown)
	}
}

func TestLoadBalanceReturnsParseErrorForInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.json")
	if err := os.WriteFile(path, []byte(`{"player":`), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	if _, err := LoadBalance(path); err == nil {
		t.Fatal("LoadBalance() error = nil, want parse error")
	}
}

func TestLoadBalanceReturnsReadErrorForMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")

	if _, err := LoadBalance(path); err == nil {
		t.Fatal("LoadBalance() error = nil, want read error")
	}
}
