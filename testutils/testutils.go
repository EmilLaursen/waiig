package testutils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func IsType[T any](t *testing.T, obj any, msgAndArgs ...any) T {
	t.Helper()
	var x T
	r, ok := obj.(T)
	var msg string
	if len(msgAndArgs) > 0 {
		fstm, ok := msgAndArgs[0].(string)
		if ok {
			msg = fmt.Sprintf(fstm, msgAndArgs[1:]...)
		}
	}
	require.True(t, ok, "type of obj=%+v is not type=%T but %s: %s", obj, x, reflect.TypeOf(obj), msg)
	return r
}
