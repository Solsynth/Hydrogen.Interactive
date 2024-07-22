package services

import (
	"github.com/pemistahl/lingua-go"
	"strings"
)

func DetectLanguage(content string) string {
	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(lingua.AllLanguages()...).
		Build()
	if lang, ok := detector.DetectLanguageOf(content); ok {
		return strings.ToLower(lang.String())
	}
	return "unknown"
}
