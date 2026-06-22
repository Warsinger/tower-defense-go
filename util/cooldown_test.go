package util

import "testing"

func TestCooldownTimerLifecycle(t *testing.T) {
	cooldown := NewCooldownTimer(3)
	if cooldown.InCooldown {
		t.Fatal("new cooldown starts in cooldown, want ready")
	}
	if got := cooldown.GetDisplay(); got != 0 {
		t.Fatalf("new cooldown display = %v, want 0", got)
	}

	cooldown.StartCooldown()
	if !cooldown.InCooldown {
		t.Fatal("StartCooldown() did not enter cooldown")
	}
	if got := cooldown.GetDisplay(); got != 3 {
		t.Fatalf("display after StartCooldown() = %v, want 3", got)
	}

	cooldown.IncrementTicker()
	cooldown.CheckCooldown()
	if !cooldown.InCooldown {
		t.Fatal("cooldown ended too early after one tick")
	}
	if got := cooldown.GetDisplay(); got != 2 {
		t.Fatalf("display after one tick = %v, want 2", got)
	}

	cooldown.IncrementTicker()
	cooldown.IncrementTicker()
	cooldown.CheckCooldown()
	if cooldown.InCooldown {
		t.Fatal("cooldown still active after reaching cooldown duration")
	}
	if got := cooldown.GetDisplay(); got != 0 {
		t.Fatalf("display after cooldown completes = %v, want 0", got)
	}
}

func TestCooldownTimerDoesNotTickWhenReady(t *testing.T) {
	cooldown := NewCooldownTimer(1)
	cooldown.IncrementTicker()
	cooldown.CheckCooldown()

	if cooldown.InCooldown {
		t.Fatal("ready cooldown became active after IncrementTicker")
	}
	if got := cooldown.GetDisplay(); got != 0 {
		t.Fatalf("ready cooldown display after IncrementTicker = %v, want 0", got)
	}
}
