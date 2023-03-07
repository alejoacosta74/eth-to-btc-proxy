package testutils

import (
	"math"
	"testing"
)

func HandleFatalError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Fatal error: %+v", err)
	}
}

// AreEqual is a helper function to compare two floats
// using an epsilon value as a margin of error
func AreEqual(got, want float64) bool {
	const epsilon = 10e-6
	delta := math.Abs(got - want)
	if want == 0 {
		return delta < epsilon
	}
	return delta/want < epsilon

}
