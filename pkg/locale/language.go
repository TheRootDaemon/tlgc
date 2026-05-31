package locale

import (
	"os"
	"strings"
)

// LanguagesFromEnv parses the LANGUAGE and LANG environment variables
// per the tldr client specification
// and appends the resulting language codes to out.
//
// https://github.com/tldr-pages/tldr/blob/main/CLIENT-SPECIFICATION.md#language
//
// If LANG is not set, the function returns without appending anything.
// LANGUAGE is split on colons and processed first,
// then LANG is appended at the end.
//
// For each entry in ll_CC format (e.g. en_US),
// both the full locale
// and the bare language are appended.
func GetLanguages(out *[]string) {
	lang := os.Getenv("LANG")
	if lang == "" {
		return
	}

	for part := range strings.SplitSeq(os.Getenv("LANGUAGE"), ":") {
		if part == "" {
			continue
		}
		appendLocale(out, part)
	}

	appendLocale(out, lang)
}

// appendLocaleVariants appends normalized locale variants derived from s.
//
// For locales in ll_CC form (for example, "en_US.UTF-8"),
// both the full locale ("en_US") and the base language ("en")
// are appended in order.
//
// Plain two-letter language codes (for example, "en")
// are appended as-is.
//
// Invalid or unsupported locale formats are ignored.
func appendLocale(out *[]string, s string) {
	if len(s) >= 5 && s[2] == '_' {
		*out = append(*out, s[:5])
		*out = append(*out, s[:2])
	} else if len(s) == 2 {
		*out = append(*out, s)
	}
}
