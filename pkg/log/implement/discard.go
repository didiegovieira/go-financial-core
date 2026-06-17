package implement

import "io"

// Discard is a convenient function to create a logger that just does nothing,
// discarding every message that it gets. Useful for testing.
func Discard() Logger {
	return discard{}
}

var _ Logger = discard{}

type discard struct{}

func (discard) Errorf(string, ...interface{}) {}
func (discard) Infof(string, ...interface{})  {}

func (d discard) Tag(string, interface{}) Logger {
	return d
}

func (d discard) Output(io.Writer) Logger {
	return d
}
