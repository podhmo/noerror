package noerr_test

import (
	"testing"

	"github.com/podhmo/noerr"
)

func TestEqual(t *testing.T) {
	add := func(x, y int) int {
		return x + y
	}

	noerr.Must(t, noerr.Equal(30).Actual(add(10, 20)))
	noerr.Should(t, noerr.Equal(30).Actual(add(10, 20)))

	noerr.Must(t, noerr.NotEqual(31).Actual(add(10, 20)))
	noerr.Should(t, noerr.NotEqual(31).Actual(add(10, 20)))
}

func TestDeepEqual(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{Name: "foo", Age: 20}
	p2 := Person{Name: "foo", Age: 20}
	noerr.Must(t, noerr.DeepEqual(p).Actual(p))
	noerr.Must(t, noerr.DeepEqual(&p).Actual(&p))
	noerr.Should(t, noerr.DeepEqual(p).Actual(p2))

	noerr.Should(t, noerr.NotDeepEqual(p).Actual(&p))
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

	noerr.Must(t, noerr.JSONEqual(p).Actual(p))
	noerr.Must(t, noerr.JSONEqual(&p).Actual(p))
	noerr.Must(t, noerr.JSONEqual(p).Actual(&p))
	noerr.Must(t, noerr.JSONEqual(&p).Actual(&p))
	noerr.Should(t, noerr.JSONEqual(p).Actual(p1))
	noerr.Should(t, noerr.JSONEqual(p).Actual(p2))

	noerr.Must(t, noerr.NotJSONEqual(nil).Actual(&p))
	noerr.Must(t, noerr.NotJSONEqual(&p).Actual(nil))
	noerr.Should(t, noerr.NotJSONEqual(p).Actual(p3))
	noerr.Should(t, noerr.NotJSONEqual(p3).Actual(p))
}
