package noerror

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

// Handy is internal object, handling compare state, not have to use this, directly
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

// ActualWithError bind actual value and no error
func (h *Handy) ActualWithError(actual interface{}, rerr error) *NG {
	if rerr != nil {
		return &NG{Name: h.Name, InnerError: rerr} // xxx
	}
	return h.Actual(actual)
}

// NG is Error value, if test is failed, wrapping state by this struct
type NG struct {
	Actual     interface{}
	Expected   interface{}
	InnerError error
	Name       string
}

// ToReport :
func (ng *NG) ToReport(toReport func(r *Reporter, err *NG, args ...interface{}) string) string {
	return toReport(DefaultReporter, ng)
}

// Describe :
func (ng *NG) Describe(name string) *NG {
	if ng == nil {
		return nil
	}
	return &NG{
		Name:       name,
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

// Log logs error, if error is occured, logs by t.Log() and return logged message by string
func Log(t testing.TB, err error, args ...interface{}) string {
	t.Helper()
	return DefaultReporter.Log(t, err, args...)
}

// Reporter :
type Reporter struct {
	ToString func(val interface{}) string
	ToReport func(r *Reporter, ng *NG, args ...interface{}) string
}

// Must not have error, if error is occured, reported by t.Fatal()
func (r *Reporter) Must(t testing.TB, err error, args ...interface{}) {
	t.Helper()
	if err == nil {
		return
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return
	}

	text, err := r.Report(err, args...)
	if err != nil {
		t.Fatalf("unexpected error, %+v", err)
	}
	t.Fatal(text)
}

// Should not have error, if error is occured, reported by t.Error()
func (r *Reporter) Should(t testing.TB, err error, args ...interface{}) {
	t.Helper()
	if err == nil {
		return
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return
	}

	text, err := r.Report(err, args...)
	if err != nil {
		t.Fatalf("unexpected error, %+v", err)
	}
	t.Error(text)
}

// Log logs error, if error is occured, logs by t.Log() and return logged message by string
func (r *Reporter) Log(t testing.TB, err error, args ...interface{}) string {
	t.Helper()
	if err == nil {
		return ""
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return ""
	}

	text, err := r.Report(err, args...)
	if err != nil {
		text, err = r.Report(err)
		if err != nil {
			panic(err)
		}
	}

	t.Log(text)
	return text
}

// Report :
func (r *Reporter) Report(err error, args ...interface{}) (string, error) {
	switch x := err.(type) {
	case *NG:
		if x.InnerError != nil {
			return "", x.InnerError
		}

		if r.ToReport != nil {
			return r.ToReport(r, x, args...), nil
		}
		return DefaultReporter.ToReport(r, x, args), nil
	default:
		return withArgs(fmt.Sprintf("unexpected error, %+v", err), args), nil
	}
}

func withArgs(text string, args []interface{}) string {
	if len(args) == 0 {
		return text
	}
	texts := []string{text, "\n"}
	for _, x := range args {
		texts = append(texts, toString(x))
	}
	return strings.Join(texts, "")
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
		ToReport: func(r *Reporter, ng *NG, args ...interface{}) string {
			name := ng.Name

			toString := r.ToString
			if toString == nil {
				toString = DefaultReporter.ToString
			}
			text := fmt.Sprintf(
				"%s, expected %s, but actual %s",
				name, toString(ng.Expected), toString(ng.Actual),
			)
			return withArgs(text, args)
		},
	}
}
