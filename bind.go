package noerror

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

// Bind : (XXX: copied value is bound)
func Bind(t *testing.T, dst interface{}) *Binder {
	t.Helper()
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr {
		Must(t, fmt.Errorf("%T is not pointer", dst))
	}
	return &Binder{t: t, ob: dst}
}

// Binder :
type Binder struct {
	t        *testing.T
	ob       interface{}
	teardown func()
}

// Actual :
func (b *Binder) Actual(ob interface{}) *Binder {
	rdst := reflect.ValueOf(b.ob)
	rsrc := reflect.ValueOf(ob)
	rdst.Elem().Set(rsrc.Elem())
	return b
}

// ActualWithTeardown :
func (b *Binder) ActualWithTeardown(ob interface{}, teardown func()) *Binder {
	b.Actual(ob).teardown = teardown
	return b
}

// ActualWithError :
func (b *Binder) ActualWithError(ob interface{}, err error) *Binder {
	b.Actual(ob)
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

	closer, ok := (b.ob).(io.Closer)
	b.t.Helper()
	if !ok {
		Must(b.t, fmt.Errorf("%T is not io.Closer", b.ob))
		return
	}
	Must(b.t, closer.Close())
}
