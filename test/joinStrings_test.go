package test

import (
	"escobra/cmd"
	"testing"
)

type test[I any, O any] struct {
	input     I
	expOutput O
}

func TestJoinStrings(t *testing.T) {
	// tests := []test{{input: } }
	t.Run("many indices", func(t *testing.T) {
		s := []string{"first-index", "second-index", "third-*"}
		expextedResult := "first-index,second-index,third-*"
		result := cmd.ParseArgsIntoSingleString(s)

		if result != expextedResult {
			t.Errorf("got:%v but expected: %v", result, expextedResult)
		}
	})

	t.Run("one index", func(t *testing.T) {
		s := []string{"first-index"}
		expResult := "first-index"
		result := cmd.ParseArgsIntoSingleString(s)
		if result != expResult {
			t.Errorf("got:%v but expected: %v", result, expResult)
		}
	})

}
