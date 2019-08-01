package handy

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestMessageFormat(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		got := Message(t, Equal(11).Actual(10))
		want := `Equal, expected 11, but actual 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
	t.Run("describe", func(t *testing.T) {
		got := Message(t, Equal(11).Actual(10).Describe("*it*"))
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

		got := Message(t, Equal(10).Actual(10).Epilog(dummy))
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

		got := Message(t, Equal(11).Actual(10).Epilog(dummy))
		want := "Equal, expected 11, but actual 10\n:bomb:"
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})

	t.Run("with ToDescription", func(t *testing.T) {
		r := &Reporter{
			ToString: DefaultReporter.ToString,
			ToDescription: func(r *Reporter, ng *NG) string {
				fmtText := "%s, want %s, but got %s"
				toString := DefaultReporter.ToString
				return fmt.Sprintf(fmtText, ng.Name, toString(ng.Expected), toString(ng.Actual))
			},
		}

		got := r.Message(t, Equal(11).Actual(10))
		want := `Equal, want 11, but got 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
}

func TestMessageWithError(t *testing.T) {
	count := func() (int, error) {
		return 0, fmt.Errorf(":bomb:")
	}

	got := Message(t, Equal(0).ActualWithNoError(count()))
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
