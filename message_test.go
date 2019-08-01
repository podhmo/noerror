package handy

import "testing"

func TestMessageFormat(t *testing.T) {
	got := Message(t, Equal(10).Expected(11))
	want := `Equal, expected 11, but actual 10`
	if got != want {
		t.Errorf("expected %q, but actual %q", want, got)
	}
}
