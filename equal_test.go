package handy_test

import (
	"testing"

	"github.com/podhmo/handy"
)

func add(x, y int) int {
	return x + y
}

func TestEqual(t *testing.T) {
	handy.Require(t, handy.Equal(add(10, 20)).Except(30))
	handy.Assert(t, handy.Equal(add(10, 20)).Except(30))

	handy.Require(t, handy.NotEqual(add(10, 20)).Except(31))
	handy.Assert(t, handy.NotEqual(add(10, 20)).Except(31))
}
