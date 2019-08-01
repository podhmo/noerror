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
	Expected interface{}
	Name     string
}

// Message :
func (ng *NG) Message(fn func(fmt string, args ...interface{}) string) string {
	var actual string
	if x, ok := ng.Actual.(fmt.Stringer); ok {
		actual = x.String()
	} else {
		actual = fmt.Sprintf("%+v", ng.Actual)
	}

	var expected string
	if x, ok := ng.Expected.(fmt.Stringer); ok {
		expected = x.String()
	} else {
		expected = fmt.Sprintf("%+v", ng.Expected)
	}
	return fn("%s, expected %s, but actual %s", ng.Name, expected, actual)
}

// Error :
func (ng *NG) Error() string {
	return ng.Message(fmt.Sprintf)
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

// Expected :
func (h *Handy) Expected(expected interface{}) error {
	ok, err := h.Compare(h.Actual, expected)
	if err != nil {
		return err // xxx
	}
	if !ok {
		return &NG{
			Actual:   h.Actual,
			Expected: expected,
			Name:     h.Name,
		}
	}
	return nil
}

// Require no error, must not be error, if error is occured, reported by t.Fatal()
func Require(t *testing.T, err error, options ...func(*Reporter)) {
	t.Helper()
	if err == nil {
		return
	}
	r := &Reporter{}
	for _, opt := range options {
		opt(r)
	}
	text := r.BuildText(err)
	t.Fatal(text)
}

// Assert no error, must not be error, if error is occured, reported by t.Error()
func Assert(t *testing.T, err error, options ...func(*Reporter)) {
	t.Helper()
	if err == nil {
		return
	}
	r := &Reporter{}
	for _, opt := range options {
		opt(r)
	}
	text := r.BuildText(err)
	t.Errorf(text)
}

// Message :
func Message(t *testing.T, err error, options ...func(*Reporter)) string {
	t.Helper()
	if err == nil {
		return ""
	}
	r := &Reporter{}
	for _, opt := range options {
		opt(r)
	}
	text := r.BuildText(err)

	t.Log(text)
	return text
}

// Reporter :
type Reporter struct {
	Message    string
	FormatText string
}

// BuildText :
func (r *Reporter) BuildText(err error) string {
	type messager interface {
		Message(fn func(s string, args ...interface{}) string) string
	}

	switch x := err.(type) {
	case messager:
		return x.Message(func(s string, args ...interface{}) string {
			name := args[0].(string) // fmt.Sprintf("%s, expected %s, but actual %s")
			if r.Message != "" {
				name = r.Message
			}
			if r.FormatText != "" {
				s = r.FormatText
			}
			return fmt.Sprintf(s, name, args[1], args[2])
		})
	case fmt.Stringer:
		return x.String()
	default:
		return x.Error()
	}
}

// WithMessage :
func WithMessage(message string) func(*Reporter) {
	return func(r *Reporter) {
		r.Message = message
	}
}

// WithFormatText : format text, default is `"%s, expected %s, but actual %s",`
func WithFormatText(fmtText string) func(*Reporter) {
	return func(r *Reporter) {
		r.FormatText = fmtText
	}
}
