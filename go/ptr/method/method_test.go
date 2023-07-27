package ptr_test

import (
	"bytes"
	"testing"

	ptr "github.com/ever0de/playground/ptr/method"
)

func TestMethodOverride(t *testing.T) {
	o := &ptr.Origin{}

	ov := ptr.NewOverride(o)
	if !bytes.Equal(ov.Get(), []byte("override")) {
		t.Fatal("override failed")
	}
}
