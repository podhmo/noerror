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
func (h *Handy) Expected(expected interface{}) *NG {
	ok, err := h.Compare(h.Actual, expected)
	if err != nil {
		return &NG{Name: h.Name, InnerError: err} // xxx
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

// NG NG
type NG struct {
	Actual     interface{}
	Expected   interface{}
	InnerError error
	Name       string
	args       []interface{}
}

// Message :
func (ng *NG) Message(buildText func(r *Reporter, err *NG) string) string {
	return buildText(DefaultReporter, ng)
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

// Epilog :
func (ng *NG) Epilog(args ...interface{}) *NG {
	if ng == nil {
		return nil
	}
	return &NG{
		Name:       ng.Name,
		args:       append(ng.args, append([]interface{}{"\n"}, args...)...),
		InnerError: ng.InnerError,
		Actual:     ng.Actual,
		Expected:   ng.Expected,
	}
}

// Error :
func (ng *NG) Error() string {
	return ng.Message(DefaultReporter.ToDescription)
}

// Require no error, must not be error, if error is occured, reported by t.Fatal()
func Require(t *testing.T, err error, options ...func(*Reporter)) {
	t.Helper()
	if err == nil {
		return
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return
	}

	r := &Reporter{}
	for _, opt := range options {
		opt(r)
	}
	text := r.BuildDescrption(err)
	t.Fatal(text)
}

// Assert no error, should not be error, if error is occured, reported by t.Error()
func Assert(t *testing.T, err error, options ...func(*Reporter)) {
	t.Helper()
	if err == nil {
		return
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return
	}

	r := &Reporter{}
	for _, opt := range options {
		opt(r)
	}
	text := r.BuildDescrption(err)
	t.Errorf(text)
}

// Message :
func Message(t *testing.T, err error, options ...func(*Reporter)) string {
	t.Helper()
	if err == nil {
		return ""
	}
	if err, ok := err.(*NG); ok && err == nil { // xxx
		return ""
	}

	r := &Reporter{}
	for _, opt := range options {
		opt(r)
	}
	text := r.BuildDescrption(err)

	t.Log(text)
	return text
}

// Reporter :
type Reporter struct {
	ToString      func(val interface{}) string
	ToDescription func(r *Reporter, ng *NG) string
}

// BuildDescrption :
func (r *Reporter) BuildDescrption(err error) string {
	switch x := err.(type) {
	case *NG:
		if x.InnerError != nil {
			return r.BuildDescrption(x.InnerError)
		}
		if r.ToDescription != nil {
			return r.ToDescription(r, x)
		}
		return DefaultReporter.ToDescription(r, x)
	case fmt.Stringer:
		return x.String()
	case error:
		return x.Error()
	default:
		panic(fmt.Sprintf("unexpected type: %T", x))
	}
}

// WithDescriptionFunction :
func WithDescriptionFunction(fn func(*Reporter, *NG) string) func(*Reporter) {
	return func(r *Reporter) {
		r.ToDescription = fn
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
		ToDescription: func(r *Reporter, ng *NG) string {
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
