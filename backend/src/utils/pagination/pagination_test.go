package pagination

import (
	"net/url"
	"testing"

	"bank-service/src/libs/i18n"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Custom values", func(t *testing.T) {
			queryString := url.Values{}
			queryString.Set("page", "2")
			queryString.Set("page_size", "10")
			offset := (2 - 1) * 10
			result, err := GetPaginationFromQuery(queryString)
			assert.Equal(t, 2, result.Page)
			assert.Equal(t, 10, result.PageSize)
			assert.Equal(t, offset, result.Offset())
			assert.NoError(t, err)
		})
		t.Run("Custom negative or zero values", func(t *testing.T) {
			queryString := url.Values{}
			queryString.Set("page", "-2")
			queryString.Set("page_size", "0")
			offset := 0 // (1 - 1) * 1
			result, err := GetPaginationFromQuery(queryString)
			assert.Equal(t, 1, result.Page)
			assert.Equal(t, 20, result.PageSize)
			assert.Equal(t, offset, result.Offset())
			assert.NoError(t, err)
		})
		t.Run("Default values", func(t *testing.T) {
			queryString := url.Values{}
			offset := 0 // (1 - 1) * 20
			result, err := GetPaginationFromQuery(queryString)
			assert.Equal(t, 1, result.Page)
			assert.Equal(t, 20, result.PageSize)
			assert.Equal(t, offset, result.Offset())
			assert.NoError(t, err)
		})
	})
	t.Run("Fail", func(t *testing.T) {
		t.Run("Page exceeds maximum", func(t *testing.T) {
			queryString := url.Values{}
			queryString.Set("page", "101")
			result, err := GetPaginationFromQuery(queryString)
			assert.Nil(t, result)
			assert.EqualError(t, err, i18n.T(i18n.Message{MessageID: "ERRORS.PAGE_TOO_LARGE"}))
		})
		t.Run("PageSize exceeds maximum", func(t *testing.T) {
			queryString := url.Values{}
			queryString.Set("page_size", "101")
			result, err := GetPaginationFromQuery(queryString)
			assert.Nil(t, result)
			assert.EqualError(t, err, i18n.T(i18n.Message{MessageID: "ERRORS.PAGE_SIZE_TOO_LARGE"}))
		})
		t.Run("Page not a number", func(t *testing.T) {
			queryString := url.Values{}
			queryString.Set("page", "not_number")
			result, err := GetPaginationFromQuery(queryString)
			assert.Nil(t, result)
			assert.Error(t, err)
		})
		t.Run("PageSize not a number", func(t *testing.T) {
			queryString := url.Values{}
			queryString.Set("page_size", "not_number")
			result, err := GetPaginationFromQuery(queryString)
			assert.Nil(t, result)
			assert.Error(t, err)
		})
	})
}
