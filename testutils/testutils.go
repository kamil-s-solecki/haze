package testutils

import (
	"bytes"
	"reflect"
	"testing"
)

func AssertEquals[T comparable](t *testing.T, got T, want T) {
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func AssertNil(t *testing.T, got any) {
	if got != nil {
		t.Errorf("got %v, wanted nil", got)
	}
}

func AssertErrorEquals(t *testing.T, got error, want string) {
	if got == nil {
		t.Errorf("got nil, expected an error '%v'", want)
	} else {
		AssertEquals(t, got.Error(), want)
	}
}

func AssertFalse(t *testing.T, got bool) {
	if got {
		t.Errorf("got true, expected false")
	}
}

func AssertTrue(t *testing.T, got bool) {
	if !got {
		t.Errorf("got false, expected true")
	}
}

func AssertMapEquals[T comparable](t *testing.T, got map[T]T, want map[T]T) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func AssertByteEquals(t *testing.T, got []byte, want []byte) {
	if !bytes.Equal(got, want) {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func AssertEmpty[T any](t *testing.T, slice []T) {
	if len(slice) != 0 {
		t.Errorf("got non empty slice")
	}
}

func AssertLen[T any](t *testing.T, slice []T, length int) {
	if len(slice) != length {
		t.Errorf("got a slice of len %v, wanted %v", len(slice), length)
	}
}

func AssertMapHasNoKey[T comparable](t *testing.T, m map[T]T, key T) {
	if _, ok := m[key]; ok {
		t.Errorf("got a map with the %v key, expected no such key", key)
	}
}
