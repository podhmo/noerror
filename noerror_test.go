package noerror_test

import (
	"testing"

	"github.com/podhmo/noerror"
)

func TestEqual(t *testing.T) {
	add := func(x, y int) int {
		return x + y
	}

	noerror.Must(t, noerror.Equal(30).Actual(add(10, 20)))
	noerror.Should(t, noerror.Equal(30).Actual(add(10, 20)))

	noerror.Must(t, noerror.NotEqual(31).Actual(add(10, 20)))
	noerror.Should(t, noerror.NotEqual(31).Actual(add(10, 20)))
}

func TestDeepEqual(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{Name: "foo", Age: 20}
	p2 := Person{Name: "foo", Age: 20}
	noerror.Must(t, noerror.DeepEqual(p).Actual(p))
	noerror.Must(t, noerror.DeepEqual(&p).Actual(&p))
	noerror.Should(t, noerror.DeepEqual(p).Actual(p2))

	noerror.Should(t, noerror.NotDeepEqual(p).Actual(&p))
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

	noerror.Must(t, noerror.JSONEqual(p).Actual(p))
	noerror.Must(t, noerror.JSONEqual(&p).Actual(p))
	noerror.Must(t, noerror.JSONEqual(p).Actual(&p))
	noerror.Must(t, noerror.JSONEqual(&p).Actual(&p))
	noerror.Should(t, noerror.JSONEqual(p).Actual(p1))
	noerror.Should(t, noerror.JSONEqual(p).Actual(p2))

	noerror.Must(t, noerror.NotJSONEqual(nil).Actual(&p))
	noerror.Must(t, noerror.NotJSONEqual(&p).Actual(nil))
	noerror.Should(t, noerror.NotJSONEqual(p).Actual(p3))
	noerror.Should(t, noerror.NotJSONEqual(p3).Actual(p))
}
