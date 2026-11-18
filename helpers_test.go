package main

import "testing"

func TestFormatPlace(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"north_carolina-usa", "North Carolina, USA"},
		{"dunedin-new_zealand", "Dunedin, New Zealand"},
		{"london-uk", "London, UK"},
		{"playa_del_carmen-mexico", "Playa Del Carmen, Mexico"},
		{"nonsense", "nonsense"},
		{"", ""},
	}

	for _, tc := range tests {
		got := formatPlace(tc.input)
		if got != tc.want {
			t.Errorf("formatPlace(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFilterArtists(t *testing.T) {
	artists := []Artist{
		{ID: 1, Name: "Queen"},
		{ID: 2, Name: "Pink Floyd"},
	}

	got := filterArtists(artists, "QUE")
	if len(got) != 1 || got[0].Name != "Queen" {
		t.Errorf("searching QUE: expected only Queen, got %v", got)
	}

	got = filterArtists(artists, "")
	if len(got) != 2 {
		t.Errorf("empty search: expected 2 artists, got %d", len(got))
	}

	got = filterArtists(artists, "zzzz")
	if len(got) != 0 {
		t.Errorf("searching zzzz: expected 0 artists, got %d", len(got))
	}
}
