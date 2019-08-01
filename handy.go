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
		"\nWhere: %s\n\tactual  : %+v\n\texpected: %+v\n",
		ng.Message,
		ng.Actual,
		ng.Excepted,
	)
}

// Equal compares by (x, y) -> x == y
func Equal(actual interface{}) *Handy {
	return &Handy{
		Name:   "Equal",
		Actual: actual,
		Compare: func(x, y interface{}) (bool, error) {
			return x == y, nil
		},
	}
}

// NotEqual compares by (x, y) -> x != y
func NotEqual(actual interface{}) *Handy {
	return &Handy{
		Name:   "NotEqual",
		Actual: actual,
		Compare: func(x, y interface{}) (bool, error) {
			return x != y, nil
		},
	}
}

// DeepEqual compares by (x, y) -> reflect.DeepEqual(x, y)
func DeepEqual(actual interface{}) *Handy {
	return &Handy{
		Name:   "DeepEqual",
		Actual: actual,
		Compare: func(x, y interface{}) (bool, error) {
			return reflect.DeepEqual(x, y), nil
		},
	}
}

// NotDeepEqual compares by (x, y) -> !reflect.DeepEqual(x, y)
func NotDeepEqual(actual interface{}) *Handy {
	return &Handy{
		Name:   "NotDeepEqual",
		Actual: actual,
		Compare: func(x, y interface{}) (bool, error) {
			return !reflect.DeepEqual(x, y), nil
		},
	}
}

// JSONEqual compares by (x, y) -> reflect.DeepEqual(normalize(x), normalize(y))
func JSONEqual(actual interface{}) *Handy {
	return &Handy{
		Name:   "JSONEqual",
		Actual: actual,
		Compare: func(x, y interface{}) (bool, error) {
			nx, err := normalize(x)
			if err != nil {
				return false, err // xxx
			}
			ny, err := normalize(y)
			if err != nil {
				return false, err // xxx
			}
			return reflect.DeepEqual(nx, ny), nil
		},
	}
}

// NotJSONEqual compares by (x, y) -> !reflect.DeepEqual(normalize(x), normalize(y))
func NotJSONEqual(actual interface{}) *Handy {
	return &Handy{
		Name:   "NotJSONEqual",
		Actual: actual,
		Compare: func(x, y interface{}) (bool, error) {
			nx, err := normalize(x)
			if err != nil {
				return false, err // xxx
			}
			ny, err := normalize(y)
			if err != nil {
				return false, err // xxx
			}
			return !reflect.DeepEqual(nx, ny), nil
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
	Name    string
	Actual  interface{}
	Compare func(x, y interface{}) (bool, error)
}

// Except :
func (h *Handy) Except(expected interface{}) error {
	ok, err := h.Compare(h.Actual, expected)
	if err != nil {
		return err // xxx
	}
	if !ok {
		return &NG{
			Actual:   h.Actual,
			Excepted: expected,
			Message:  h.Name,
		}
	}
	return nil
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
