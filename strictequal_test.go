package handy_test

import (
	"testing"

	"github.com/podhmo/handy"
)

func add(x, y int) int {
	return x + y
}

func TestStrictEqual(t *testing.T) {
	handy.Require(t, handy.StrictEqual(add(10, 20)).Except(30))
	handy.Assert(t, handy.StrictEqual(add(10, 20)).Except(30))
}
