package noerror

import (
	"testing"
)

type Closer struct {
	err    error
	closed bool
}

func (c *Closer) Close() error {
	c.closed = true
	return c.err
}

func TestBind(t *testing.T) {
	t.Run("ok, with closer", func(t *testing.T) {
		var got *Closer
		defer func() {
			if !got.closed {
				t.Errorf("expect closed, but not closed")
			}
		}()

		setup := func() *Closer {
			got = &Closer{err: nil}
			return got
		}

		var ob interface{}
		defer Bind(t, &ob).Actual(setup()).Teardown()
		if _, ok := ob.(*Closer); !ok {
			t.Fatalf("expect bound, but not bound: %T", ob)
		}
	})

	t.Run("ok, with error", func(t *testing.T) {
		var got *Closer
		defer func() {
			if !got.closed {
				t.Errorf("expect closed, but not closed")
			}
		}()

		setup := func() (*Closer, error) {
			got = &Closer{err: nil}
			return got, nil
		}

		var ob interface{}
		defer Bind(t, &ob).ActualWithError(setup()).Teardown()
		if _, ok := ob.(*Closer); !ok {
			t.Fatalf("expect bound, but not bound: %T", ob)
		}
	})

	t.Run("ok, with teardown", func(t *testing.T) {
		var got *Closer
		defer func() {
			if !got.closed {
				t.Errorf("expect closed, but not closed")
			}
		}()

		setup := func() (*Closer, func()) {
			got = &Closer{err: nil}
			return got, func() { Must(t, got.Close()) }
		}

		var ob interface{}
		defer Bind(t, &ob).ActualWithTeardown(setup()).Teardown()
		if _, ok := ob.(*Closer); !ok {
			t.Fatalf("expect bound, but not bound: %T", ob)
		}
	})
	t.Run("ok, with error, with teardown", func(t *testing.T) {
		var got *Closer
		defer func() {
			if !got.closed {
				t.Errorf("expect closed, but not closed")
			}
		}()

		setup := func() (*Closer, error, func()) {
			got = &Closer{err: nil}
			return got, nil, func() { Must(t, got.Close()) }
		}

		var ob interface{}
		defer Bind(t, &ob).ActualWithErrorWithTeardown(setup()).Teardown()
		if _, ok := ob.(*Closer); !ok {
			t.Fatalf("expect bound, but not bound: %T", ob)
		}
	})
}
