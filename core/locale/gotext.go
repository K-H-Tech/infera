package locale

import (
	"context"
	"strings"

	"github.com/leonelquinteros/gotext"
)

var locales map[string]*gotext.Locale

// Init loads all locales at startup
func Init() {
	locales = map[string]*gotext.Locale{}
	for _, lang := range []string{"en_US", "fa_IR"} {
		loc := gotext.NewLocale("locales", lang)
		loc.AddDomain("messages")
		locales[lang] = loc
	}
}

type localeKey struct{}

// WithLocale attaches a *gotext.Locale to the context
func WithLocale(ctx context.Context, loc *gotext.Locale) context.Context {
	return context.WithValue(ctx, localeKey{}, loc)
}

// FromContext retrieves the *gotext.Locale from context
func FromContext(ctx context.Context) *gotext.Locale {
	if loc, ok := ctx.Value(localeKey{}).(*gotext.Locale); ok {
		return loc
	}
	return locales["en_US"]
}

// FromAcceptLang picks locale based on Accept-Language header
func FromAcceptLang(header string) *gotext.Locale {
	if header == "" {
		return locales["en_US"]
	}

	lang := strings.Split(header, ",")[0]
	lang = strings.ReplaceAll(lang, "-", "_")

	switch lang {
	case "fa", "fa_IR":
		return locales["fa_IR"]
	case "en", "en_US":
		return locales["en_US"]
	default:
		return locales["en_US"]
	}
}
