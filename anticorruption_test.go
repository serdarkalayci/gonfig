package gonfig

import "testing"

// A workaround to solve the mismatch between assert package's expectation and go 1.18's testing package
// assert package expects the signature from Errorf method to be Errorf(format string, args ...interface{})
// but with the introduction of generics in go 1.18 it changed into Errorf(format string, args ...any)
// testingACL struct satisfies the interface asert package is expecting and calls the underlying testing.T.Errorf method
type testingACL struct {
	t *testing.T
}

func (tacl testingACL) Errorf(format string, args ...interface{}) {
	tacl.t.Errorf(format, args...)
}
