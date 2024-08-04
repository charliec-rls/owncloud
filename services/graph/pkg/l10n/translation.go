package l10n

import (
	"embed"

	l10npkg "github.com/owncloud/ocis/v2/ocis-pkg/l10n"
)

var (
	//go:embed locale
	_localeFS embed.FS
)

const (
	// subfolder where the translation files are stored
	_localeSubPath = "locale"

	// domain of the graph service (transifex)
	_domain = "graph"
)

func Translate(content, locale, defaultLocale string) string {
	t := l10npkg.NewTranslatorFromCommonConfig(defaultLocale, _domain, "", _localeFS, _localeSubPath)
	return t.Translate(content, locale)
}

func NewTranslateLocation(locale, defaultLocale string) func(string, ...any) string {
	t := l10npkg.NewTranslatorFromCommonConfig(defaultLocale, _domain, "", _localeFS, _localeSubPath)
	return t.Locale(locale).Get
}