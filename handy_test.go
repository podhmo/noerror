package handy_test

import (
	"testing"

	"github.com/podhmo/handy"
)

func TestEqual(t *testing.T) {
	add := func(x, y int) int {
		return x + y
	}

	handy.Require(t, handy.Equal(30).Actual(add(10, 20)))
	handy.Assert(t, handy.Equal(30).Actual(add(10, 20)))

	handy.Require(t, handy.NotEqual(31).Actual(add(10, 20)))
	handy.Assert(t, handy.NotEqual(31).Actual(add(10, 20)))
}

func TestDeepEqual(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{Name: "foo", Age: 20}
	p2 := Person{Name: "foo", Age: 20}
	handy.Require(t, handy.DeepEqual(p).Actual(p))
	handy.Require(t, handy.DeepEqual(&p).Actual(&p))
	handy.Assert(t, handy.DeepEqual(p).Actual(p2))

	handy.Assert(t, handy.NotDeepEqual(p).Actual(&p))
}

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

	handy.Require(t, handy.JSONEqual(p).Actual(p))
	handy.Require(t, handy.JSONEqual(&p).Actual(p))
	handy.Require(t, handy.JSONEqual(p).Actual(&p))
	handy.Require(t, handy.JSONEqual(&p).Actual(&p))
	handy.Assert(t, handy.JSONEqual(p).Actual(p1))
	handy.Assert(t, handy.JSONEqual(p).Actual(p2))

	handy.Require(t, handy.NotJSONEqual(nil).Actual(&p))
	handy.Require(t, handy.NotJSONEqual(&p).Actual(nil))
	handy.Assert(t, handy.NotJSONEqual(p).Actual(p3))
	handy.Assert(t, handy.NotJSONEqual(p3).Actual(p))
}
