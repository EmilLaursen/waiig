package testutils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func IsType[T any](t *testing.T, obj any) T {
	t.Helper()
	var x T
	r, ok := obj.(T)
	require.True(t, ok, "type of obj=%+v is not type=%T but %s", obj, x, reflect.TypeOf(obj))
	return r
}
