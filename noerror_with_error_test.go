package noerror_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/podhmo/noerror"
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
func (t *fakeTB) Fatal(val ...interface{}) {
	t.called = append(t.called, "fatal")
}
func (t *fakeTB) Error(val ...interface{}) {
	t.called = append(t.called, "error")
}

func TestEqualWithError(t *testing.T) {
	count := func() (int, error) {
		return 0, fmt.Errorf(":bomb:")
	}

	t.Run("require, on error, call t.Fatalf()", func(t *testing.T) {
		ft := &fakeTB{}
		noerror.Must(ft, noerror.Equal(0).ActualWithError(count()))
		if len(ft.called) == 0 {
			t.Fatal("testing.TB's method is must be called, but not called")
		}
		if !strings.HasPrefix(ft.called[0], "fatal") {
			t.Errorf("*testing.T.Fatal() or *testing.T.Fatalf() is must be called, but %v", ft.called)
		}
	})

	t.Run("even use Should(), on error, call t.Fatalf()", func(t *testing.T) {
		ft := &fakeTB{}
		noerror.Should(ft, noerror.Equal(0).ActualWithError(count()))
		if len(ft.called) == 0 {
			t.Fatal("testing.TB's method is must be called, but not called")
		}
		if !strings.HasPrefix(ft.called[0], "fatal") {
			t.Errorf("*testing.T.Fatal() or *testing.T.Fatalf() is must be called, but %v", ft.called)
		}
	})
}
