package money

import "strings"

// SymbolPosition controls where the currency symbol is placed relative to the
// formatted number.
type SymbolPosition uint8

const (
	// SymbolBefore places the symbol immediately before the number: "$1,234.56"
	SymbolBefore SymbolPosition = iota
	// SymbolBeforeSpace places the symbol before the number with a space: "$ 1,234.56"
	SymbolBeforeSpace
	// SymbolAfter places the symbol immediately after the number: "1,234.56€"
	SymbolAfter
	// SymbolAfterSpace places the symbol after the number with a space: "1,234.56 €"
	SymbolAfterSpace
)

// LocaleConfig describes number-formatting rules for a locale.
type LocaleConfig struct {
	Decimal   string         // decimal separator ("." or ",")
	Thousand  string         // thousands separator (",", ".", "\u00a0", "'", …)
	SymbolPos SymbolPosition // symbol placement relative to the number
}

// LocaleFormats maps IETF BCP 47 locale tags to formatting rules for ~60 of
// the most common web/commerce locales. Tags are matched exactly; if no exact
// match is found, the language subtag is tried (e.g. "de" for "de-DE").
//
// You may add custom entries before calling LocaleFormat.
var LocaleFormats = map[string]LocaleConfig{
	// English-speaking regions
	"en":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"en-US": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"en-GB": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"en-AU": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"en-CA": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"en-NZ": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"en-SG": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"en-ZA": {Decimal: ".", Thousand: "\u00a0", SymbolPos: SymbolBefore},
	"en-IN": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"en-HK": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	// German
	"de":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"de-DE": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"de-AT": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBeforeSpace},
	"de-CH": {Decimal: ".", Thousand: "'", SymbolPos: SymbolBeforeSpace},
	"de-LU": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	// French
	"fr":    {Decimal: ",", Thousand: "\u202f", SymbolPos: SymbolAfterSpace},
	"fr-FR": {Decimal: ",", Thousand: "\u202f", SymbolPos: SymbolAfterSpace},
	"fr-BE": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"fr-CA": {Decimal: ",", Thousand: "\u202f", SymbolPos: SymbolBeforeSpace},
	"fr-CH": {Decimal: ".", Thousand: "'", SymbolPos: SymbolAfterSpace},
	"fr-LU": {Decimal: ",", Thousand: "\u202f", SymbolPos: SymbolAfterSpace},
	// Spanish
	"es":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"es-ES": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"es-MX": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"es-AR": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBefore},
	"es-CO": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBefore},
	"es-CL": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBefore},
	"es-PE": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	// Italian
	"it":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"it-IT": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"it-CH": {Decimal: ".", Thousand: "'", SymbolPos: SymbolBeforeSpace},
	// Dutch
	"nl":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolBeforeSpace},
	"nl-NL": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBeforeSpace},
	"nl-BE": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBeforeSpace},
	// Portuguese
	"pt":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolBeforeSpace},
	"pt-BR": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBeforeSpace},
	"pt-PT": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	// Russian / Ukrainian
	"ru":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"ru-RU": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"uk":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"uk-UA": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	// Turkish
	"tr":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolBefore},
	"tr-TR": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBefore},
	// Polish
	"pl":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"pl-PL": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	// Scandinavian
	"sv":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"sv-SE": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"nb":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"nb-NO": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"no-NO": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"da":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"da-DK": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"fi":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"fi-FI": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	// Japanese
	"ja":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"ja-JP": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	// Chinese
	"zh":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"zh-CN": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"zh-TW": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"zh-HK": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	// Korean
	"ko":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"ko-KR": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	// Arabic
	"ar":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolAfterSpace},
	"ar-SA": {Decimal: ".", Thousand: ",", SymbolPos: SymbolAfterSpace},
	"ar-AE": {Decimal: ".", Thousand: ",", SymbolPos: SymbolAfterSpace},
	"ar-EG": {Decimal: ".", Thousand: ",", SymbolPos: SymbolAfterSpace},
	// South / South-East Asia
	"hi":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"hi-IN": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"th":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"th-TH": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"id":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolBefore},
	"id-ID": {Decimal: ",", Thousand: ".", SymbolPos: SymbolBefore},
	"ms":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"ms-MY": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"vi":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"vi-VN": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"bn":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"bn-BD": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	// Middle East / Africa
	"he":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"he-IL": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"sw":    {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	"sw-KE": {Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore},
	// Central / Eastern Europe
	"el":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"el-GR": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"ro":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"ro-RO": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"hu":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"hu-HU": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"cs":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"cs-CZ": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"sk":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"sk-SK": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"bg":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"bg-BG": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"hr":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"hr-HR": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"sr":    {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"sr-RS": {Decimal: ",", Thousand: ".", SymbolPos: SymbolAfterSpace},
	"kk":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"kk-KZ": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"uz":    {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
	"uz-UZ": {Decimal: ",", Thousand: "\u00a0", SymbolPos: SymbolAfterSpace},
}

// localeTemplate converts a SymbolPosition to the template string used by
// Formatter ("$1", "$ 1", "1$", or "1 $").
func localeTemplate(pos SymbolPosition) string {
	switch pos {
	case SymbolBeforeSpace:
		return "$ 1"
	case SymbolAfter:
		return "1$"
	case SymbolAfterSpace:
		return "1 $"
	default: // SymbolBefore
		return "$1"
	}
}

// lookupLocale returns the LocaleConfig for a BCP 47 locale tag. It first
// tries an exact match, then falls back to the language subtag only.
func lookupLocale(locale string) (LocaleConfig, bool) {
	if cfg, ok := LocaleFormats[locale]; ok {
		return cfg, true
	}
	if idx := strings.IndexByte(locale, '-'); idx > 0 {
		if cfg, ok := LocaleFormats[locale[:idx]]; ok {
			return cfg, true
		}
	}
	return LocaleConfig{}, false
}

// LocaleFormat formats the monetary value according to the given IETF BCP 47
// locale tag (e.g. "en-US", "de-DE", "fr-FR"). It uses the currency's own
// symbol but applies locale-specific decimal/thousands separators and symbol
// placement. Falls back to Display() for unrecognised locales.
//
//	money.New(123456, "EUR").LocaleFormat("en-US") // "€1,234.56"
//	money.New(123456, "EUR").LocaleFormat("de-DE") // "1.234,56 €"
//	money.New(123456, "EUR").LocaleFormat("fr-FR") // "1 234,56 €"
func (m *Money) LocaleFormat(locale string) string {
	cfg, ok := lookupLocale(locale)
	if !ok {
		return m.Display()
	}
	c := m.currency.get()
	f := NewFormatter(c.Fraction, cfg.Decimal, cfg.Thousand, c.Grapheme, localeTemplate(cfg.SymbolPos))
	return f.Format(m.amount.val)
}

// AccountingFormat formats the monetary value using locale-specific separators
// and renders negative values in parentheses — the standard accounting
// convention used in financial reports and spreadsheets.
//
//	money.New(-123456, "USD").AccountingFormat("en-US") // "($1,234.56)"
//	money.New(123456,  "USD").AccountingFormat("en-US") // "$1,234.56"
//	money.New(-123456, "EUR").AccountingFormat("de-DE") // "(1.234,56 €)"
func (m *Money) AccountingFormat(locale string) string {
	cfg, ok := lookupLocale(locale)
	if !ok {
		cfg = LocaleConfig{Decimal: ".", Thousand: ",", SymbolPos: SymbolBefore}
	}
	c := m.currency.get()
	f := NewFormatter(c.Fraction, cfg.Decimal, cfg.Thousand, c.Grapheme, localeTemplate(cfg.SymbolPos))
	if m.amount.val < 0 {
		return "(" + f.Format(-m.amount.val) + ")"
	}
	return f.Format(m.amount.val)
}
