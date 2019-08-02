package noerror

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestLogFormat(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		got := Log(t, Equal(11).Actual(10))
		want := `Equal, expected 11, but actual 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
	t.Run("describe", func(t *testing.T) {
		got := Log(t, Equal(11).Actual(10).Describe("*it*"))
		want := `*it*, expected 11, but actual 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})

	t.Run("with epilog, skip, when not error", func(t *testing.T) {
		dummy := &fstringer{
			toString: func() string {
				var b strings.Builder
				var r io.Reader
				// panic
				if _, err := io.Copy(&b, r); err != nil {
					panic(err)
				}
				return b.String()
			},
		}

		got := Log(t, Equal(10).Actual(10), dummy)
		want := ``
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
	t.Run("describe with argument, fmt.Stringer", func(t *testing.T) {
		dummy := &fstringer{
			toString: func() string {
				return ":bomb:"
			},
		}

		got := Log(t, Equal(11).Actual(10), dummy)
		want := "Equal, expected 11, but actual 10\n:bomb:"
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})

	t.Run("with ToReport", func(t *testing.T) {
		r := &Reporter{
			ToString: DefaultReporter.ToString,
			ToReport: func(r *Reporter, ng *NG, args ...interface{}) string {
				fmtText := "%s, want %s, but got %s"
				toString := DefaultReporter.ToString
				return fmt.Sprintf(fmtText, ng.Name, toString(ng.Expected), toString(ng.Actual))
			},
		}

		got := r.Log(t, Equal(11).Actual(10))
		want := `Equal, want 11, but got 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})

	t.Run("raw error", func(t *testing.T) {
		got := Log(t, &fstringer{toString: func() string { return "*raw*" }})
		want := "unexpected error, !!*raw*"
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})

	t.Run("pkg/errors", func(t *testing.T) {
		raw := &fstringer{toString: func() string { return "*raw*" }}
		got := Log(t, errors.WithMessage(raw, "WRAP"))
		want := "unexpected error, !!*raw*\nWRAP"
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
}

func TestLogWithNoError(t *testing.T) {
	count := func() (int, error) {
		return 0, fmt.Errorf(":bomb:")
	}

	got := Log(t, Equal(0).ActualWithNoError(count()))
	want := "unexpected error, :bomb:"
	if got != want {
		t.Errorf("expected %q, but actual %q", want, got)
	}
}

type fstringer struct {
	toString func() string
}

func (x *fstringer) String() string {
	return x.toString()
}

func (x *fstringer) Error() string {
	return "!!" + x.toString()
}
