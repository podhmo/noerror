package handy_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/podhmo/handy"
)

type fakeTB struct {
	testing.TB
	called []string
}

func (t *fakeTB) Helper() {
}
func (t *fakeTB) Fatalf(s string, val ...interface{}) {
	t.called = append(t.called, "fatalf")
}
func (t *fakeTB) Error(val ...interface{}) {
	t.called = append(t.called, "error")
}

func TestEqualWithNoError(t *testing.T) {
	count := func() (int, error) {
		return 0, fmt.Errorf(":bomb:")
	}

	t.Run("on, error, call t.Fatalf()", func(t *testing.T) {
		ft := &fakeTB{}
		handy.Assert(ft, handy.EqualWithNoError(count()).Expected(0))
		if len(ft.called) == 0 {
			t.Fatal("testing.TB's method is must be called, but not called")
		}
		if !strings.HasPrefix(ft.called[0], "fatal") {
			t.Errorf("*testing.T.Fatal() or *testing.T.Fatalf() is must be called, but %v", ft.called)
		}
	})
}
