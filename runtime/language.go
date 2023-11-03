package runtime

import "golang.org/x/text/language"

var matcher language.Matcher

func init() {
	supported := []language.Tag{
		language.English,
		language.AmericanEnglish,
		language.Russian,
		language.Finnish,
	}

	matcher = language.NewMatcher(supported)
}

func GetLanguage(lang string) string {
	userLang := language.MustParse(lang)
	matchedLang, _, _ := matcher.Match(userLang)
	return matchedLang.String()
}
