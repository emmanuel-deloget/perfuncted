package input

import (
	"strings"
	"testing"
)

func TestXkbKeysym(t *testing.T) {
	tests := []struct {
		r    rune
		want string
	}{
		{'A', "U0041"},
		{'z', "U007A"},
		{'0', "U0030"},
		{'€', "U20AC"},
	}
	for _, tc := range tests {
		got := xkbKeysym(tc.r)
		if got != tc.want {
			t.Errorf("xkbKeysym(%q) = %q, want %q", tc.r, got, tc.want)
		}
	}
}

func TestNamedKey(t *testing.T) {
	tests := []struct {
		name string
		kc   uint32
		sym  string
		ok   bool
	}{
		{"shift", kcShift, "Shift_L", true},
		{"Ctrl", kcCtrl, "Control_L", true},
		{"RETURN", kcReturn, "Return", true},
		{"enter", kcReturn, "Return", true},
		{"f1", kcF1, "F1", true},
		{"f12", kcF1 + 11, "F12", true},
		{"nonexistent", 0, "", false},
	}
	for _, tc := range tests {
		kc, sym, ok := namedKey(tc.name)
		if ok != tc.ok || kc != tc.kc || sym != tc.sym {
			t.Errorf("namedKey(%q) = (%d, %q, %v), want (%d, %q, %v)",
				tc.name, kc, sym, ok, tc.kc, tc.sym, tc.ok)
		}
	}
}

func TestModBit(t *testing.T) {
	tests := []struct {
		key  string
		want uint32
	}{
		{"shift", modShift},
		{"ctrl", modControl},
		{"alt", modMod1},
		{"super", modMod4},
		{"f1", 0},
	}
	for _, tc := range tests {
		got := modBit(tc.key)
		if got != tc.want {
			t.Errorf("modBit(%q) = %d, want %d", tc.key, got, tc.want)
		}
	}
}

func TestUniqueRunes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "helo"},
		{"aabbcc", "abc"},
		{"", ""},
		{"abcabc", "abc"},
		{"日本日", "日本"},
	}
	for _, tc := range tests {
		got := uniqueRunes(tc.input)
		gotStr := string(got)
		if gotStr != tc.want {
			t.Errorf("uniqueRunes(%q) = %q, want %q", tc.input, gotStr, tc.want)
		}
	}
}

func TestXkbBuildContainsModifiers(t *testing.T) {
	km := xkbModsOnly()
	if !strings.Contains(km, "Shift_L") {
		t.Error("keymap missing Shift_L")
	}
	if !strings.Contains(km, "Control_L") {
		t.Error("keymap missing Control_L")
	}
	if !strings.Contains(km, "Alt_L") {
		t.Error("keymap missing Alt_L")
	}
	if !strings.Contains(km, "Super_L") {
		t.Error("keymap missing Super_L")
	}
}

func TestXkbWithRunes(t *testing.T) {
	km := xkbWithRunes([]rune{'A', 'b'})
	if !strings.Contains(km, "U0041") {
		t.Error("keymap missing U0041 for 'A'")
	}
	if !strings.Contains(km, "U0062") {
		t.Error("keymap missing U0062 for 'b'")
	}
}

func TestXkbWithNamed(t *testing.T) {
	km := xkbWithNamed(kcReturn, "Return")
	if !strings.Contains(km, "Return") {
		t.Error("keymap missing Return")
	}
	if !strings.Contains(km, "KNAM") {
		t.Error("keymap missing KNAM keycode")
	}

	// Modifier keycodes should just return mods-only.
	km2 := xkbWithNamed(kcShift, "Shift_L")
	if strings.Contains(km2, "KNAM") {
		t.Error("modifier keycode should not add KNAM")
	}
}
