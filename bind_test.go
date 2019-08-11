package noerror

import (
	"testing"
)

type Closer struct {
	err     error
	closed  bool
	created bool
}

func (c *Closer) Close() error {
	c.closed = true
	return c.err
}

func TestBind(t *testing.T) {
	t.Run("ok, with closer", func(t *testing.T) {
		var got Closer
		defer func() {
			if !got.closed {
				t.Errorf("expect closed, but not closed")
			}
		}()

		setup := func() *Closer {
			return &Closer{created: true, err: nil}
		}
		defer Bind(t, &got).Actual(setup()).Teardown()

		if !got.created {
			t.Errorf("expect created, but not created")
		}
	})

	t.Run("ok, with error", func(t *testing.T) {
		var got Closer
		defer func() {
			if !got.closed {
				t.Errorf("expect closed, but not closed")
			}
		}()

		setup := func() (*Closer, error) {
			return &Closer{created: true, err: nil}, nil
		}

		defer Bind(t, &got).ActualWithError(setup()).Teardown()

		if !got.created {
			t.Errorf("expect created, but not created")
		}
	})

	t.Run("ok, with teardown", func(t *testing.T) {
		var got Closer
		defer func() {
			if !got.closed {
				t.Errorf("expect closed, but not closed")
			}
		}()

		setup := func() (*Closer, func()) {
			got := &Closer{created: true, err: nil}
			return got, func() { Must(t, got.Close()) }
		}

		defer Bind(t, &got).ActualWithTeardown(setup()).Teardown()

		if !got.created {
			t.Errorf("expect created, but not created")
		}
	})
	t.Run("ok, with error, with teardown", func(t *testing.T) {
		var got Closer
		defer func() {
			if !got.closed {
				t.Errorf("expect closed, but not closed")
			}
		}()

		setup := func() (*Closer, error, func()) {
			got := &Closer{created: true, err: nil}
			return got, nil, func() { Must(t, got.Close()) }
		}

		defer Bind(t, &got).ActualWithErrorWithTeardown(setup()).Teardown()
		if !got.created {
			t.Errorf("expect created, but not created")
		}
	})
}
