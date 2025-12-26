package tests

import "testing"

func assert(t *testing.T, condition bool, message string, args ...any) {
	if condition {
		t.Errorf(message, args...)
	}
}
