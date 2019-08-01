package handy

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestMessageFormat(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		got := Message(t, Equal(10).Expected(11))
		want := `Equal, expected 11, but actual 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
	t.Run("describe", func(t *testing.T) {
		got := Message(t, Equal(10).Expected(11).Describe("*it*"))
		want := `*it*, expected 11, but actual 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
	t.Run("with ToDescription", func(t *testing.T) {
		got := Message(t,
			Equal(10).Expected(11),
			WithDescriptionFunction(func(r *Reporter, ng *NG) string {
				fmtText := "%s, want %s, but got %s"
				toString := DefaultReporter.ToString
				return fmt.Sprintf(fmtText, ng.Name, toString(ng.Expected), toString(ng.Actual))
			}),
		)
		want := `Equal, want 11, but got 10`
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

		got := Message(t, Equal(10).Expected(10).Epilog(dummy))
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

		got := Message(t, Equal(10).Expected(11).Epilog(dummy))
		want := "Equal, expected 11, but actual 10\n:bomb:"
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
}

func TestMessageWithError(t *testing.T) {
	count := func() (int, error) {
		return 0, fmt.Errorf(":bomb:")
	}

	got := Message(t, EqualWithNoError(count()).Expected(0))
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
