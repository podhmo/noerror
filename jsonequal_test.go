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
	type Person3 struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	p := Person{Name: "foo", Age: 20}
	p1 := Person{Name: "foo", Age: 20}
	p2 := Person2{Name: "foo", Age: 20}
	p3 := Person3{Name: "foo", Age: 20}

	handy.Require(t, handy.JSONEqual(p).Except(p))
	handy.Require(t, handy.JSONEqual(&p).Except(p))
	handy.Require(t, handy.JSONEqual(p).Except(&p))
	handy.Require(t, handy.JSONEqual(&p).Except(&p))
	handy.Assert(t, handy.JSONEqual(p).Except(p1))
	handy.Assert(t, handy.JSONEqual(p).Except(p2))

	handy.Require(t, handy.JSONNotEqual(nil).Except(&p))
	handy.Require(t, handy.JSONNotEqual(&p).Except(nil))
	handy.Assert(t, handy.JSONNotEqual(p).Except(p3))
	handy.Assert(t, handy.JSONNotEqual(p3).Except(p))
}
