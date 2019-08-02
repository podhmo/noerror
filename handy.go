package handy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var (
	DefaultReporter *Reporter
)

// Equal compares by (x, y) -> x == y
func Equal(expected interface{}) *Handy {
	return &Handy{
		Name:     "Equal",
		Expected: expected,
		Compare: func(x, y interface{}) (bool, error) {
			return x == y, nil
		},
	}
}

// NotEqual compares by (x, y) -> x != y
func NotEqual(expected interface{}) *Handy {
	return &Handy{
		Name:     "NotEqual",
		Expected: expected,
		Compare: func(x, y interface{}) (bool, error) {
			return x != y, nil
		},
	}
}

// DeepEqual compares by (x, y) -> reflect.DeepEqual(x, y)
func DeepEqual(expected interface{}) *Handy {
	return &Handy{
		Name:     "DeepEqual",
		Expected: expected,
		Compare: func(x, y interface{}) (bool, error) {
			return reflect.DeepEqual(x, y), nil
		},
	}
}

// NotDeepEqual compares by (x, y) -> !reflect.DeepEqual(x, y)
func NotDeepEqual(expected interface{}) *Handy {
	return &Handy{
		Name:     "NotDeepEqual",
		Expected: expected,
		Compare: func(x, y interface{}) (bool, error) {
			return !reflect.DeepEqual(x, y), nil
		},
	}
}

// JSONEqual compares by (x, y) -> reflect.DeepEqual(normalize(x), normalize(y))
func JSONEqual(expected interface{}) *Handy {
	return &Handy{
		Name:     "JSONEqual",
		Expected: expected,
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
func NotJSONEqual(expected interface{}) *Handy {
	return &Handy{
		Name:     "NotJSONEqual",
		Expected: expected,
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
	Name     string
	Expected interface{}
	Compare  func(x, y interface{}) (bool, error)
}

// Actual binds actual value
func (h *Handy) Actual(actual interface{}) *NG {
	ok, err := h.Compare(h.Expected, actual)
	if err != nil {
		return &NG{Name: h.Name, InnerError: err} // xxx
	}
	if !ok {
		return &NG{
			Expected: h.Expected,
			Actual:   actual,
			Name:     h.Name,
		}
	}
	return nil
}

// ActualWithNoError bind actual value and no error
func (h *Handy) ActualWithNoError(actual interface{}, rerr error) *NG {
	if rerr != nil {
		return &NG{Name: h.Name, InnerError: rerr} // xxx
	}
	return h.Actual(actual)
}

// NG NG
type NG struct {
	Actual     interface{}
	Expected   interface{}
	InnerError error
	Name       string
	args       []interface{}
}

// ToReport :
func (ng *NG) ToReport(toReport func(r *Reporter, err *NG) string) string {
	return toReport(DefaultReporter, ng)
}

// Describe :
func (ng *NG) Describe(name string) *NG {
	if ng == nil {
		return nil
	}
	return &NG{
		Name:       name,
		args:       ng.args,
		InnerError: ng.InnerError,
		Actual:     ng.Actual,
		Expected:   ng.Expected,
	}
}

// Error :
func (ng *NG) Error() string {
	return ng.ToReport(DefaultReporter.ToReport)
}

// Must not have error, if error is occured, reported by t.Fatal()
func Must(t testing.TB, err error, args ...interface{}) {
	t.Helper()
	DefaultReporter.Must(t, err, args...)
}

// Should not have error, if error is occured, reported by t.Error()
func Should(t testing.TB, err error, args ...interface{}) {
	t.Helper()
	DefaultReporter.Should(t, err, args...)
}

// Log :
func Log(t testing.TB, err error, args ...interface{}) string {
	t.Helper()
	return DefaultReporter.Log(t, err, args...)
}

// Reporter :
type Reporter struct {
	ToString      func(val interface{}) string
	ToReport func(r *Reporter, ng *NG) string
}

// Must no error, must not be error, if error is occured, reported by t.Fatal()
func (r *Reporter) Must(t testing.TB, err error, args ...interface{}) {
	t.Helper()
	if err == nil {
		return
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return
	}

	text, err := r.Descrption(err, args...)
	if err != nil {
		t.Fatalf("unexpected error, %+v", err)
	}
	t.Fatal(text)
}

// Should no error, should not be error, if error is occured, reported by t.Error()
func (r *Reporter) Should(t testing.TB, err error, args ...interface{}) {
	t.Helper()
	if err == nil {
		return
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return
	}

	text, err := r.Descrption(err, args...)
	if err != nil {
		t.Fatalf("unexpected error, %+v", err)
	}
	t.Error(text)
}

// Log :
func (r *Reporter) Log(t testing.TB, err error, args ...interface{}) string {
	t.Helper()
	if err == nil {
		return ""
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return ""
	}

	text, err := r.Descrption(err, args...)
	if err != nil {
		text, err = r.Descrption(err)
		if err != nil {
			panic(err)
		}
		text = fmt.Sprintf("unexpected error, %+v", text)
	}

	t.Log(text)
	return text
}

// Descrption :
func (r *Reporter) Descrption(err error, args ...interface{}) (string, error) {
	switch x := err.(type) {
	case *NG:
		if x.InnerError != nil {
			return "", x.InnerError
		}

		if len(args) > 0 {
			x.args = append(append(x.args, "\n"), args...)
		}

		if r.ToReport != nil {
			return r.ToReport(r, x), nil
		}
		return DefaultReporter.ToReport(r, x), nil
	case fmt.Stringer:
		return x.String(), nil
	case error:
		return x.Error(), nil
	default:
		panic(fmt.Sprintf("unexpected type: %T", x))
	}
}

func toString(val interface{}) string {
	if x, ok := val.(fmt.Stringer); ok {
		return x.String()
	}
	return fmt.Sprintf("%+v", val)
}

func init() {
	DefaultReporter = &Reporter{
		ToString: toString,
		ToReport: func(r *Reporter, ng *NG) string {
			name := ng.Name

			toString := r.ToString
			if toString == nil {
				toString = DefaultReporter.ToString
			}
			fmtText := "%s, expected %s, but actual %s"
			description := fmt.Sprintf(fmtText, name, toString(ng.Expected), toString(ng.Actual))
			if ng.args == nil {
				return description
			}
			texts := []string{description}
			for _, x := range ng.args {
				texts = append(texts, toString(x))
			}
			return strings.Join(texts, "")
		},
	}
}
