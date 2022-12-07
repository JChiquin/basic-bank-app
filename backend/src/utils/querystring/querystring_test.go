package querystring

import (
	"bank-service/src/utils/constant"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
)

func TestNewQueryStringDecoder(t *testing.T) {
	//Fixtures
	name := "Jorge"
	date := time.Now()
	type queryStringDTO struct {
		Name string    `schema:"name"`
		Date time.Time `schema:"date"`
	}
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Dates with time", func(t *testing.T) {
			testCases := []struct {
				TestName string
				Values   map[string][]string
			}{
				{
					TestName: "Time without quotation marks",
					Values: map[string][]string{
						"name": {name},
						"date": {date.Format(time.RFC3339)},
					},
				},
				{
					TestName: "Time with quotation marks",
					Values: map[string][]string{
						"name": {name},
						"date": {fmt.Sprintf(`""%s""`, date.Format(time.RFC3339))},
					},
				},
			}
			for _, tC := range testCases {
				t.Run(tC.TestName, func(t *testing.T) {
					//Fixtures
					queryStringObj := queryStringDTO{}

					//Action
					decoder := newQueryStringDecoder()
					err := decoder.Decode(&queryStringObj, tC.Values)

					//Data Assertion
					assert.NoError(t, err)
					assert.Equal(t, name, queryStringObj.Name)
					assert.WithinDuration(t, date, queryStringObj.Date, 1*time.Second)
				})
			}
		})
		t.Run("Timeless date", func(t *testing.T) {
			//Fixtures
			queryStringObj := queryStringDTO{}
			dateTimeless := map[string][]string{
				"name": {name},
				"date": {date.Format(constant.FormatYYYYzMMzDD)},
			}

			//Action
			decoder := newQueryStringDecoder()
			err := decoder.Decode(&queryStringObj, dateTimeless)

			//Data Assertion
			assert.NoError(t, err)
			assert.Equal(t, name, queryStringObj.Name)
			assert.Equal(t, now.With(date).BeginningOfDay(), queryStringObj.Date)
		})
	})
	t.Run("Should fail on invalid time", func(t *testing.T) {
		//Fixtures
		queryStringObj := queryStringDTO{}
		queryStringValues := map[string][]string{
			"name": {name},
			"date": {"invalid time"},
		}

		//Action
		decoder := newQueryStringDecoder()
		err := decoder.Decode(&queryStringObj, queryStringValues)

		//Data Assertion
		assert.Error(t, err)
	})
}

func TestDecode(t *testing.T) {
	type Body struct {
		Name string `json:"name" schema:"name,required"`
		Age  int    `json:"age"`
	}
	validName := "name"
	validAgeStr := "123"
	validAge := 123
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Full struct", func(t *testing.T) {
			//Fixture
			values := url.Values{}
			values.Add("name", validName)
			values.Add("age", validAgeStr)
			body := Body{}

			// Action
			err := Decode(&body, values)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, validName, body.Name)
			assert.Equal(t, validAge, body.Age)
		})
		t.Run("Without not required field", func(t *testing.T) {
			//Fixture
			values := url.Values{}
			values.Add("name", validName)
			body := Body{}

			// Action
			err := Decode(&body, values)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, validName, body.Name)
			assert.Zero(t, body.Age)
		})
	})
	t.Run("Should fail on", func(t *testing.T) {
		t.Run("Parsing error", func(t *testing.T) {
			//Fixture
			values := url.Values{}
			values.Add("name", validName)
			values.Add("age", "NotNumber")
			body := Body{}

			// Action
			err := Decode(&body, values)

			// Assert
			assert.Error(t, err)
		})
		t.Run("Missing required value", func(t *testing.T) {
			//Fixture
			values := url.Values{}
			values.Add("age", validAgeStr)
			body := Body{}

			// Action
			err := Decode(&body, values)

			// Assert
			assert.Error(t, err)
		})
	})
}
