package services

import (
	"strings"

	"github.com/pemistahl/lingua-go"
)

var detector lingua.LanguageDetector

func CreateLanguageDetector() lingua.LanguageDetector {
	return lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		WithLowAccuracyMode().
		Build()
}

func DetectLanguage(content string) string {
	if detector == nil {
		detector = CreateLanguageDetector()
	}

	if lang, ok := detector.DetectLanguageOf(content); ok {
		return strings.ToLower(lang.String())
	}
	return "unknown"
}
