package test

import (
	"fmt"
	"testing"
)

func TestQueryBuilders(t *testing.T) {
	//

	type L struct {
		x []string
	}

	l := L{}

	for _, i := range l.x {
		fmt.Println(i)
	}

}
