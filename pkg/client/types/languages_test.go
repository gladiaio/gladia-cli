package types

import (
	"strings"
	"testing"
)

func TestParseLanguages_empty(t *testing.T) {
	langs, err := ParseLanguages("")
	if err != nil {
		t.Fatal(err)
	}
	if langs != nil {
		t.Fatalf("got %v, want nil", langs)
	}
}

func TestParseLanguages_single(t *testing.T) {
	langs, err := ParseLanguages("en")
	if err != nil {
		t.Fatal(err)
	}
	if len(langs) != 1 || langs[0] != LanguageEn {
		t.Fatalf("got %v", langs)
	}
}

func TestParseLanguages_multiple(t *testing.T) {
	langs, err := ParseLanguages(" en , FR , de ")
	if err != nil {
		t.Fatal(err)
	}
	if len(langs) != 3 {
		t.Fatalf("got %v", langs)
	}
	if langs[0] != LanguageEn || langs[1] != LanguageFr || langs[2] != LanguageDe {
		t.Fatalf("got %v", langs)
	}
}

func TestParseLanguages_dedupes(t *testing.T) {
	langs, err := ParseLanguages("en,en,fr")
	if err != nil {
		t.Fatal(err)
	}
	if len(langs) != 2 {
		t.Fatalf("got %v", langs)
	}
}

func TestParseLanguages_invalid(t *testing.T) {
	_, err := ParseLanguages("english")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown language") {
		t.Fatalf("got %v", err)
	}
}

func TestParseLanguages_onlyCommas(t *testing.T) {
	_, err := ParseLanguages(",,,")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLanguage_single(t *testing.T) {
	lang, err := ParseLanguage("fr")
	if err != nil {
		t.Fatal(err)
	}
	if lang != LanguageFr {
		t.Fatalf("got %v", lang)
	}
}

func TestParseLanguage_multipleRejected(t *testing.T) {
	_, err := ParseLanguage("en,fr")
	if err == nil {
		t.Fatal("expected error")
	}
}
