package normalize

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const placeholder = "\x00"

func Nombre(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "ñ", placeholder)
	s = strings.ReplaceAll(s, "Ñ", placeholder)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)

	s = strings.ReplaceAll(s, placeholder, "ñ")

	return s
}

func Tags(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "ñ", placeholder)
	s = strings.ReplaceAll(s, "Ñ", placeholder)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)

	s = strings.ReplaceAll(s, placeholder, "ñ")
	return s
}

func Busqueda(s string) string {
	return Nombre(s)
}

func Iguales(a, b string) bool {
	return Nombre(a) == Nombre(b)
}
