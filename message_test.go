package handy

import "testing"

func TestMessageFormat(t *testing.T) {
	t.Run("equal error", func(t *testing.T) {
		got := Message(t, Equal(10).Expected(11))
		want := `Equal, expected 11, but actual 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
	t.Run("with message", func(t *testing.T) {
		got := Message(t,
			Equal(10).Expected(11),
			WithMessage("*with message*"),
		)
		want := `*with message*, expected 11, but actual 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
	t.Run("with formatText", func(t *testing.T) {
		got := Message(t,
			Equal(10).Expected(11),
			WithFormatText("%s, want %s, but got %s"),
		)
		want := `Equal, want 11, but got 10`
		if got != want {
			t.Errorf("expected %q, but actual %q", want, got)
		}
	})
}
