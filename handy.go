package handy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// NG NG
type NG struct {
	Actual   interface{}
	Excepted interface{}
	Message  string
}

// Error :
func (ng *NG) Error() string {
	return fmt.Sprintf(
		"%s:\n\tActual:%+v\n\tExpected:%+v\n",
		ng.Message,
		ng.Actual,
		ng.Excepted,
	)
}

// StrictEqual compares by (x, y) -> x == y
func StrictEqual(actual interface{}) *Handy {
	return &Handy{
		actual: actual,
		Compare: func(x, y interface{}) error {
			if x == y {
				return nil
			}
			return &NG{
				Actual:   x,
				Excepted: y,
				Message:  "StrictEqual",
			}
		},
	}
}

// StrictNotEqual compares by (x, y) -> x != y
func StrictNotEqual(actual interface{}) *Handy {
	return &Handy{
		actual: actual,
		Compare: func(x, y interface{}) error {
			if x != y {
				return nil
			}
			return &NG{
				Actual:   x,
				Excepted: y,
				Message:  "StrictNotEqual",
			}
		},
	}
}

// DeepEqual compares by (x, y) -> reflect.DeepEqual(x, y)
func DeepEqual(actual interface{}) *Handy {
	return &Handy{
		actual: actual,
		Compare: func(x, y interface{}) error {
			if reflect.DeepEqual(x, y) {
				return nil
			}
			return &NG{
				Actual:   x,
				Excepted: y,
				Message:  "DeepEqual",
			}
		},
	}
}

// DeepNotEqual compares by (x, y) -> !reflect.DeepEqual(x, y)
func DeepNotEqual(actual interface{}) *Handy {
	return &Handy{
		actual: actual,
		Compare: func(x, y interface{}) error {
			if !reflect.DeepEqual(x, y) {
				return nil
			}
			return &NG{
				Actual:   x,
				Excepted: y,
				Message:  "DeepNotEqual",
			}
		},
	}
}

// JSONEqual compares by (x, y) -> reflect.Equal(normalize(x), normalize(y))
func JSONEqual(actual interface{}) *Handy {
	return &Handy{
		actual: actual,
		Compare: func(x, y interface{}) error {
			nx, err := normalize(x)
			if err != nil {
				return err // xxx
			}
			ny, err := normalize(y)
			if err != nil {
				return err // xxx
			}
			if reflect.DeepEqual(nx, ny) {
				return nil
			}
			return &NG{
				Actual:   x,
				Excepted: y,
				Message:  "JSONEqual",
			}
		},
	}
}

// JSONEqual compares by (x, y) -> reflect.Equal(normalize(x), normalize(y))

func JSONNotEqual(actual interface{}) *Handy {
	return &Handy{
		actual: actual,
		Compare: func(x, y interface{}) error {
			nx, err := normalize(x)
			if err != nil {
				return err // xxx
			}
			ny, err := normalize(y)
			if err != nil {
				return err // xxx
			}
			if !reflect.DeepEqual(nx, ny) {
				return nil
			}
			return &NG{
				Actual:   x,
				Excepted: y,
				Message:  "JSONNotEqual",
			}
		},
	}
}

func normalize(src interface{}) (interface{}, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	var dst interface{}
	if err := json.Unmarshal(b, &dst); err != nil {
		return nil, err
	}
	return dst, nil
}

// Handy :
type Handy struct {
	Compare func(x, y interface{}) error
	actual  interface{}
}

// Except :
func (h *Handy) Except(expected interface{}) error {
	return h.Compare(h.actual, expected)
}

// Require no error, must not be error, if error is occured, reported by t.Fatal()
func Require(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	t.Fatalf("%s", err)
}

// Assert no error, must not be error, if error is occured, reported by t.Error()
func Assert(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	t.Errorf("%s", err)
}

// Report :
func Report(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		return ""
	}
	t.Logf("%s", err)
	return fmt.Sprintf("%s", err)
}
