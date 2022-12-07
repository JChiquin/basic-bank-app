package i18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

const (
	languageEs string = "es"
	languageEn string = "en"
)

func TestSetupI18n(t *testing.T) {
	assert.NotPanics(t, SetupI18n)
	assert.NotNil(t, loc)
	assert.NotNil(t, bundle)
	assert.Len(t, bundle.LanguageTags(), 2)
}

func TestSetLanguage(t *testing.T) {
	previousLoc := loc

	assert.NotPanics(t, func() { SetLanguage(languageEs) })
	assert.NotNil(t, loc)
	assert.NotEqual(t, loc, previousLoc)
}

func TestT(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		testCases := []struct {
			Language    language.Tag
			Description string
		}{
			{
				Language:    language.Spanish,
				Description: "Mensaje",
			},
			{
				Language:    language.English,
				Description: "Message",
			},
		}
		for _, tC := range testCases {
			t.Run(tC.Language.String(), func(t *testing.T) {
				ID := "FOO"
				SetLanguage(tC.Language.String())
				bundle.MustAddMessages(tC.Language, &i18n.Message{ID: ID, Other: tC.Description})

				//Action
				got := T(Message{MessageID: ID})

				//Data Assertion
				assert.Equal(t, tC.Description, got)
			})
		}
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("Message not found", func(t *testing.T) {
			//Action
			got := T(Message{MessageID: "ID not found"})

			//Data Assertion
			assert.Contains(t, got, "not found")
		})
	})
}
