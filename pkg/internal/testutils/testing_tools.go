package testutils

import (
	"testing"
)

func HandleFatalError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Fatal error: %+v", err)
	}
}
