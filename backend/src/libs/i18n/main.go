package i18n

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle
var loc *i18n.Localizer

//Message is to be used instead of i18n.LocalizeConfig
type Message i18n.LocalizeConfig

func init() {
	SetupI18n()
}

/*
SetupI18n initializes the bundle with default language (spanish), loads the files for each language
and initializes the localizer without language (it uses the default)
*/
func SetupI18n() {
	bundle = i18n.NewBundle(language.Spanish)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("/app/src/libs/i18n/en.json")
	bundle.MustLoadMessageFile("/app/src/libs/i18n/es.json")
	loc = i18n.NewLocalizer(bundle)
}

//SetLanguage set the language provides by header, per requests
func SetLanguage(language string) {
	loc = i18n.NewLocalizer(bundle, language)
}

//T receives a message and translates it according language selected
func T(message Message) string {
	localizeConfig := i18n.LocalizeConfig(message)
	localizedMessage, err := loc.Localize(&localizeConfig)
	if err != nil {
		return err.Error()
	}
	return localizedMessage
}
