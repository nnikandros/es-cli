package test

import (
	"escobra/cmd"
	"testing"
)

func TestIndexNameValid(t *testing.T) {
	t.Run("", func(t *testing.T) {
		result := cmd.IsIndexNameValid(".monitor")
		expResult := false

		if result != expResult {
			t.Errorf("got %v but expected %v", result, expResult)
		}
	})
	t.Run("", func(t *testing.T) {
		result := cmd.IsIndexNameValid("connector-")
		expResult := false

		if result != expResult {
			t.Errorf("got %v but expected %v", result, expResult)
		}
	})

}
