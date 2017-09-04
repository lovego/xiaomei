package logger

import (
	"testing"
)

func TestPrintf(t *testing.T) {
	log := New(`test: `, nil, nil)
	log.Printf("Printf: %s", `error`)
}
