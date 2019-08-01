package handy_test

import (
	"testing"

	"github.com/podhmo/handy"
)

func TestDeepEqual(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{Name: "foo", Age: 20}
	p2 := Person{Name: "foo", Age: 20}
	handy.Require(t, handy.DeepEqual(p).Except(p))
	handy.Require(t, handy.DeepEqual(&p).Except(&p))
	handy.Assert(t, handy.DeepEqual(p).Except(p2))

	handy.Assert(t, handy.DeepNotEqual(p).Except(&p))
}
