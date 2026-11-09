package main

import (
	"testing"

	"github.com/mhsanaei/3x-ui/v3/database"
	"github.com/mhsanaei/3x-ui/v3/web/service"
)

func setupMainTestDB(t *testing.T) {
	t.Helper()
	t.Setenv("XUI_DB_TYPE", "")
	t.Setenv("XUI_DB_FOLDER", t.TempDir())
	t.Cleanup(func() {
		if err := database.CloseDB(); err != nil {
			t.Fatal(err)
		}
	})
}

// `x-ui setting -subPort N` used to be parsed and then silently ignored,
// so install.sh's "2096 is taken, move to 2097" never took effect.
func TestUpdateSettingSubPort(t *testing.T) {
	setupMainTestDB(t)
	s := service.SettingService{}

	if err := updateSetting(0, 2097, "", "", "", "", false); err != nil {
		t.Fatalf("updateSetting: %v", err)
	}
	got, err := s.GetSubPort()
	if err != nil {
		t.Fatal(err)
	}
	if got != 2097 {
		t.Fatalf("subPort = %d, want 2097", got)
	}

	// Not passing -subPort must leave it alone.
	if err := updateSetting(12345, 0, "", "", "", "", false); err != nil {
		t.Fatalf("updateSetting: %v", err)
	}
	if got, _ := s.GetSubPort(); got != 2097 {
		t.Fatalf("subPort changed to %d when flag not given", got)
	}
	if p, _ := s.GetPort(); p != 12345 {
		t.Fatalf("port = %d, want 12345", p)
	}
}

func TestUpdateSettingSubPortRejectsOutOfRange(t *testing.T) {
	setupMainTestDB(t)
	s := service.SettingService{}

	if err := updateSetting(0, 2097, "", "", "", "", false); err != nil {
		t.Fatalf("updateSetting: %v", err)
	}
	for _, bad := range []int{-1, 65536, 70000} {
		if err := updateSetting(0, bad, "", "", "", "", false); err == nil {
			t.Fatalf("subPort %d: expected error", bad)
		}
	}
	if got, _ := s.GetSubPort(); got != 2097 {
		t.Fatalf("subPort = %d after rejected values, want 2097", got)
	}
}
