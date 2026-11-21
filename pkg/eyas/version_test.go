package eyas

import (
	"testing"
)

func TestParseLongVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected LongVersion
	}{
		{"v1.0.0-7-ge622ad89", LongVersion{Version: "v1.0.0", Commit: "7-ge622ad89"}},
		{"v1.0.0", LongVersion{Version: "v1.0.0", Commit: ""}},
		{"1.2.3-4-gabcdef", LongVersion{Version: "1.2.3", Commit: "4-gabcdef"}},
		{"v0.10.0-rc1-39-g4f62f61", LongVersion{Version: "v0.10.0-rc1", Commit: "39-g4f62f61"}},
		{"1.2.3", LongVersion{Version: "1.2.3", Commit: ""}},
		{"1.2.3-rc2", LongVersion{Version: "1.2.3-rc2", Commit: ""}},
	}

	for _, test := range tests {
		result := ParseLongVersion(test.input)
		if result.Version != test.expected.Version || result.Commit != test.expected.Commit {
			t.Errorf("ParseLongVersion(%q) = %+v; want %+v", test.input, result, test.expected)
		}
	}
}

func TestParseBuildInfo(t *testing.T) {
	sample := `build time: 2025-07-23 09:32:30
prod: v0.10.0-11-ge0e5df0
eyas: v1.0.0-3-g3abd273
nays: v0.10.0-12-gf0b4a07
nte: v0.10.0-26-g748dc8a
web-ui: v0.10.0-rc1-39-g4f62f61
xsnos-cli: v0.9.0-0-g856307c
`
	v := ParseBuildInfo([]byte(sample))
	if v.BuildTime != "2025-07-23 09:32:30" {
		t.Fatalf("expected build time '2025-07-23 09:32:30', got '%s'", v.BuildTime)
	}
	if v.Prod.Version != "v0.10.0" || v.Prod.Commit != "11-ge0e5df0" {
		t.Fatalf("expected prod 'v0.10.0-11-ge0e5df0', got '%s'", v.Prod)
	}
	if v.Eyas.Version != "v1.0.0" || v.Eyas.Commit != "3-g3abd273" {
		t.Fatalf("expected eyas 'v1.0.0-3-g3abd273', got '%s'", v.Eyas)
	}
	if len(v.Components) != 4 {
		t.Fatalf("expected 4 components, got %d", len(v.Components))
	}

	expectedComponents := []struct {
		name    string
		version LongVersion
	}{
		{"nays", LongVersion{Version: "v0.10.0", Commit: "12-gf0b4a07"}},
		{"nte", LongVersion{Version: "v0.10.0", Commit: "26-g748dc8a"}},
		{"web-ui", LongVersion{Version: "v0.10.0-rc1", Commit: "39-g4f62f61"}},
		{"xsnos-cli", LongVersion{Version: "v0.9.0", Commit: "0-g856307c"}},
	}

	for i, exp := range expectedComponents {
		if v.Components[i].Name != exp.name || v.Components[i].Version != exp.version {
			t.Fatalf("component %d mismatch: expected %+v, got %+v", i, exp, v.Components[i])
		}
	}
}
