package game

import (
	"fmt"
)

type I18nEntry struct {
	fr string
	en string
}

var I18n = map[string]I18nEntry{
	"l1.e1.1": {fr: "Bonjour je m'appelle Caniche à maman", en: "Hello my name is Caniche mother fucker"},
}

func I18nGet(k string) string {
	if res, ok := I18n[k]; ok {
		switch gm.locale {
		case "fr":
			return res.fr
		case "en":
			return res.en
		default:
			return "I18n undefined locale"
		}
	} else {
		return fmt.Sprintf("I18n key %s is not defined", k)
	}
}
