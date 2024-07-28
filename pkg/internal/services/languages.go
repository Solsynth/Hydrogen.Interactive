package services

import (
	"strings"

	"github.com/pemistahl/lingua-go"
)

func DetectLanguage(content string) string {
	return "unknown"

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(lingua.AllLanguages()...).
		Build()
	if lang, ok := detector.DetectLanguageOf(content); ok {
		return strings.ToLower(lang.String())
	}
	return "unknown"
}
