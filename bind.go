package noerror

import (
	"fmt"
	"io"
	"testing"
)

// Bind :
func Bind(t *testing.T, dst *interface{}) *Binder {
	return &Binder{t: t, ob: dst}
}

// Binder :
type Binder struct {
	t        *testing.T
	ob       *interface{}
	teardown func()
}

// Actual :
func (b *Binder) Actual(ob interface{}) *Binder {
	*b.ob = ob
	return b
}

// ActualWithTeardown :
func (b *Binder) ActualWithTeardown(ob interface{}, teardown func()) *Binder {
	b.Actual(ob).teardown = teardown
	return b
}

// ActualWithError :
func (b *Binder) ActualWithError(ob interface{}, err error) *Binder {
	*b.ob = ob
	b.t.Helper()
	Must(b.t, err)
	return b
}

// ActualWithErrorWithTeardown :
func (b *Binder) ActualWithErrorWithTeardown(ob interface{}, err error, teardown func()) *Binder {
	b.ActualWithError(ob, err).teardown = teardown
	return b
}

// Teardown :
func (b *Binder) Teardown() {
	if b.teardown != nil {
		b.teardown()
		return
	}

	closer, ok := (*b.ob).(io.Closer)
	b.t.Helper()
	if !ok {
		Must(b.t, fmt.Errorf("%T is not io.Closer", *b.ob))
		return
	}
	Must(b.t, closer.Close())
}
