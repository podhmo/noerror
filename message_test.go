package handy

import (
	"fmt"
	"testing"
)

func TestMessageFormat(t *testing.T) {
	t.Run("equal error", func(t *testing.T) {
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
}
