package handy_test

import (
	"testing"

	"github.com/podhmo/handy"
)

func TestJSONEqual(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	type Person2 struct {
		Name string
		Age  int
	}

	p := Person{Name: "foo", Age: 20}
	p2 := Person{Name: "foo", Age: 20}
	p3 := Person2{Name: "foo", Age: 20}

	handy.Require(t, handy.JSONEqual(p).Except(p))
	handy.Require(t, handy.JSONEqual(&p).Except(p))
	handy.Require(t, handy.JSONEqual(p).Except(&p))
	handy.Require(t, handy.JSONEqual(&p).Except(&p))
	handy.Assert(t, handy.JSONEqual(p).Except(p2))
	handy.Assert(t, handy.JSONEqual(p).Except(p3))
}
