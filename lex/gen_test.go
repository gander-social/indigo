package lex

import "testing"

func TestParsePackages(t *testing.T) {
	text := `[{"package": "gndr", "prefix": "app.gndr", "outdir": "api/gndr", "import": "github.com/gander-social/gander-indigo-sovereign/api/gndr"}]`
	parsed, err := ParsePackages([]byte(text))
	if err != nil {
		t.Fatalf("error parsing json: %s", err)
	}
	if len(parsed) != 1 {
		t.Fatalf("expected 1, got %d", len(parsed))
	}
	expected := Package{"gndr", "app.gndr", "api/gndr", "github.com/gander-social/gander-indigo-sovereign/api/gndr"}
	if expected != parsed[0] {
		t.Fatalf("expected %#v, got %#v", expected, parsed[0])
	}

}
