package achievement

import "testing"

func TestGetAchievementDefByID(t *testing.T) {
	got := GetAchievementDefByID("not_rubbed_yet")
	if got == nil {
		t.Fatal("GetAchievementDefByID('not_rubbed_yet') = nil, want non-nil")
	}
	if got.ID != "not_rubbed_yet" {
		t.Errorf("got.ID = %q, want %q", got.ID, "not_rubbed_yet")
	}
	if got.Name != "AchievementNameNotRubbedYet" {
		t.Errorf("got.Name = %q, want %q", got.Name, "AchievementNameNotRubbedYet")
	}

	gotLast := GetAchievementDefByID("annihilator_cannon")
	if gotLast == nil {
		t.Fatal("GetAchievementDefByID('annihilator_cannon') = nil, want non-nil")
	}

	gotMissing := GetAchievementDefByID("nonexistent")
	if gotMissing != nil {
		t.Errorf("GetAchievementDefByID('nonexistent') = %v, want nil", gotMissing)
	}
}

func TestAllAchievementsNotEmpty(t *testing.T) {
	if len(AllAchievements) == 0 {
		t.Error("AllAchievements is empty")
	}

	ids := make(map[string]bool)
	for _, ach := range AllAchievements {
		if ach.ID == "" {
			t.Error("found achievement with empty ID")
		}
		if ids[ach.ID] {
			t.Errorf("duplicate achievement ID: %q", ach.ID)
		}
		ids[ach.ID] = true
	}
}
