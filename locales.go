package main

import (
	"github.com/omnea/faker"
	"github.com/omnea/faker/locales"
	"strings"
)

func setLocale(country string) {
	switch strings.ToLower(country) {
	case "de":
		faker.Locale = locales.De
	case "es":
		faker.Locale = locales.Es
	case "en":
		faker.Locale = locales.En
	default:
		faker.Locale = locales.En
	}
}
