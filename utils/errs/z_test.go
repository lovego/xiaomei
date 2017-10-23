package errs

import (
	"testing"
)

func TestData(t *testing.T) {
	err := New(`code`, `message`)
	if err.Data() != nil {
		t.Error(`Data() should return nil.`)
	}
}
