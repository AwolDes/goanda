package goanda

import "testing"

func logTestResult(t *testing.T, name string) {
    if t.Failed() {
        t.Logf("\n❌ Test failed: %s", name)
    } else {
        t.Logf("\n✅ Test passed: %s", name)
	}
}