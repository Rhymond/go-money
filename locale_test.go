package money

import "testing"

func TestMoney_LocaleFormat(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		locale   string
		expected string
	}{
		// en-US: symbol before, dot decimal, comma thousand
		{123456, "EUR", "en-US", "€1,234.56"},
		{100, "USD", "en-US", "$1.00"},
		{-123456, "USD", "en-US", "-$1,234.56"},
		// de-DE: symbol after space, comma decimal, dot thousand
		{123456, "EUR", "de-DE", "1.234,56 €"},
		{-100, "EUR", "de-DE", "-1,00 €"},
		// language-subtag fallback: "de-XX" not in map, falls back to "de"
		{123456, "EUR", "de-XX", "1.234,56 €"},
		// unknown locale → fallback to Display()
		{123456, "EUR", "xx-XX", "€1,234.56"},
		// zero
		{0, "USD", "en-US", "$0.00"},
		// JPY (fraction=0)
		{1000, "JPY", "ja-JP", "¥1,000"},
		// SymbolBeforeSpace (de-AT)
		{12345, "EUR", "de-AT", "€ 123,45"},
		// es-AR: SymbolBefore, comma decimal, dot thousand
		{100, "ARS", "es-AR", "$1,00"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		got := m.LocaleFormat(tc.locale)
		if got != tc.expected {
			t.Errorf("LocaleFormat(%d %s, %q): expected %q got %q",
				tc.amount, tc.code, tc.locale, tc.expected, got)
		}
	}
}

func TestMoney_LocaleFormat_NBSP(t *testing.T) {
	// fr-FR uses narrow no-break space (U+202F) — just check it returns non-empty
	m := New(123456, "EUR")
	got := m.LocaleFormat("fr-FR")
	if got == "" {
		t.Error("LocaleFormat(fr-FR): unexpected empty string")
	}
}

func TestMoney_LocaleFormat_AllPositions(t *testing.T) {
	// Ensure all four SymbolPosition variants produce the expected output.

	// SymbolBefore (en-US): "$1.00"
	m := New(100, "USD")
	if got := m.LocaleFormat("en-US"); got != "$1.00" {
		t.Errorf("SymbolBefore: expected $1.00 got %s", got)
	}

	// SymbolBeforeSpace (de-AT): "€ 1,00"
	eur := New(100, "EUR")
	if got := eur.LocaleFormat("de-AT"); got != "€ 1,00" {
		t.Errorf("SymbolBeforeSpace (de-AT): expected '€ 1,00' got %q", got)
	}

	// SymbolAfterSpace (de-DE): "1,00 €"
	if got := eur.LocaleFormat("de-DE"); got != "1,00 €" {
		t.Errorf("SymbolAfterSpace (de-DE): expected '1,00 €' got %q", got)
	}

	// SymbolAfter — no standard locale uses this, so insert a synthetic entry.
	LocaleFormats["_test-SYMAFTER"] = LocaleConfig{
		Decimal: ".", Thousand: ",", SymbolPos: SymbolAfter,
	}
	defer delete(LocaleFormats, "_test-SYMAFTER")
	m2 := New(123456, "USD")
	if got := m2.LocaleFormat("_test-SYMAFTER"); got != "1,234.56$" {
		t.Errorf("SymbolAfter: expected '1,234.56$' got %q", got)
	}
}

func TestMoney_AccountingFormat(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		locale   string
		expected string
	}{
		// positive — no parens
		{123456, "USD", "en-US", "$1,234.56"},
		// negative — wrapped in parens
		{-123456, "USD", "en-US", "($1,234.56)"},
		// zero — no parens
		{0, "USD", "en-US", "$0.00"},
		// European negative
		{-100, "EUR", "de-DE", "(1,00 €)"},
		// unknown locale falls back to default config (SymbolBefore, ".", ",")
		{-100, "USD", "xx-ZZ", "($1.00)"},
		{100, "USD", "xx-ZZ", "$1.00"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		got := m.AccountingFormat(tc.locale)
		if got != tc.expected {
			t.Errorf("AccountingFormat(%d %s, %q): expected %q got %q",
				tc.amount, tc.code, tc.locale, tc.expected, got)
		}
	}
}

func TestLocaleTemplate(t *testing.T) {
	tcs := []struct {
		pos      SymbolPosition
		expected string
	}{
		{SymbolBefore, "$1"},
		{SymbolBeforeSpace, "$ 1"},
		{SymbolAfter, "1$"},
		{SymbolAfterSpace, "1 $"},
		{SymbolPosition(99), "$1"}, // default
	}
	for _, tc := range tcs {
		got := localeTemplate(tc.pos)
		if got != tc.expected {
			t.Errorf("localeTemplate(%d): expected %q got %q", tc.pos, tc.expected, got)
		}
	}
}

func TestLookupLocale(t *testing.T) {
	// Exact match
	cfg, ok := lookupLocale("en-US")
	if !ok || cfg.Decimal != "." {
		t.Errorf("lookupLocale(en-US): expected found with decimal='.' got ok=%v cfg=%+v", ok, cfg)
	}
	// Language fallback: "de-XX" not in map, falls back to "de"
	cfg, ok = lookupLocale("de-XX")
	if !ok || cfg.Decimal != "," {
		t.Errorf("lookupLocale(de-XX): expected language fallback to de got ok=%v cfg=%+v", ok, cfg)
	}
	// Unknown
	_, ok = lookupLocale("zz-ZZ")
	if ok {
		t.Error("lookupLocale(zz-ZZ): expected not found")
	}
	// No hyphen (exact only, no fallback possible)
	_, ok = lookupLocale("zzz")
	if ok {
		t.Error("lookupLocale(zzz): expected not found")
	}
}

func TestMoney_Sign(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int
	}{
		{100, 1},
		{-100, -1},
		{0, 0},
	}
	for _, tc := range tcs {
		m := New(tc.amount, "USD")
		if got := m.Sign(); got != tc.expected {
			t.Errorf("Sign(%d): expected %d got %d", tc.amount, tc.expected, got)
		}
	}
}
