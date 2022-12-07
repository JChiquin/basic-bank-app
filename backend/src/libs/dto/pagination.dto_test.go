package dto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	t.Run("Offset", func(t *testing.T) {
		testCases := []struct {
			Input    *Pagination
			Expected int
		}{
			{
				Input:    NewPagination(1, 20, 0), //Don't care about third arg (TotalCount)
				Expected: 0,
			},
			{
				Input:    NewPagination(2, 20, 0),
				Expected: 20,
			},
			{
				Input:    NewPagination(3, 20, 0),
				Expected: 40,
			},
			{
				Input:    NewPagination(1, 5, 0),
				Expected: 0,
			},
			{
				Input:    NewPagination(2, 5, 0),
				Expected: 5,
			},
		}
		for _, tC := range testCases {
			t.Run(fmt.Sprintf("With Page=%d and PageSize=%d, Offset should be %d",
				tC.Input.Page, tC.Input.PageSize, tC.Expected), func(t *testing.T) {
				got := tC.Input.Offset()
				assert.Equal(t, tC.Expected, got)
			})
		}
	})
	t.Run("PageCount", func(t *testing.T) {
		testCases := []struct {
			Input    *Pagination
			Expected int
		}{
			{
				Input:    NewPagination(0, 20, 10), //Don't care about first arg (Page)
				Expected: 1,
			},
			{
				Input:    NewPagination(0, 20, 30),
				Expected: 2,
			},
			{
				Input:    NewPagination(0, 20, 39),
				Expected: 2,
			},
			{
				Input:    NewPagination(0, 20, 40),
				Expected: 2,
			},
			{
				Input:    NewPagination(0, 20, 41),
				Expected: 3,
			},
			{
				Input:    NewPagination(0, 5, 100),
				Expected: 20,
			},
		}
		for _, tC := range testCases {
			t.Run(fmt.Sprintf("With TotalCount=%d and PageSize=%d, PageCount should be %d",
				tC.Input.TotalCount, tC.Input.PageSize, tC.Expected), func(t *testing.T) {
				got := tC.Input.PageCount()
				assert.Equal(t, tC.Expected, got)
			})
		}
	})
}
