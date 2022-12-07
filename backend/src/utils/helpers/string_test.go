package helpers

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	//Fixture
	list := []string{"Foo", "Bar", "Foo Bar", "Mock", "Testing", "With spaces"}
	type arg struct {
		ToFind string
		List   []string
	}
	t.Run("Should success on", func(t *testing.T) {

		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomIndex := r1.Intn(len(list)) //[0, len)

		testCases := []struct {
			TestName string
			Input    arg
		}{
			{
				TestName: "First word",
				Input: arg{
					ToFind: list[0],
					List:   list,
				},
			},
			{
				TestName: "Last word",
				Input: arg{
					ToFind: list[len(list)-1],
					List:   list,
				},
			},
			{
				TestName: "Random word",
				Input: arg{
					ToFind: list[randomIndex],
					List:   list,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.TestName, func(t *testing.T) {
				found := StringInSlice(tC.Input.ToFind, tC.Input.List)
				assert.True(t, found)
			})
		}
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Input    arg
		}{
			{
				TestName: "Empty list",
				Input: arg{
					ToFind: list[0],
					List:   []string{},
				},
			},
			{
				TestName: "Empty string to find",
				Input: arg{
					ToFind: "",
					List:   list,
				},
			},
			{
				TestName: "Not found",
				Input: arg{
					ToFind: "Test not found",
					List:   list,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.TestName, func(t *testing.T) {
				found := StringInSlice(tC.Input.ToFind, tC.Input.List)
				assert.False(t, found)
			})
		}
	})
}

func TestContextKey(t *testing.T) {
	contextKey := "key1"
	toString := ContextKey(contextKey)

	assert.Equal(t, contextKey, toString.String())
}
