// Package log is just a small configuration layer upon other loggers that
// provide the real log features.
package implement

import (
	"fmt"
	"strings"
)

type Level uint32

const (
	LevelError Level = iota
	LevelInfo
)

var allLevels = []Level{
	LevelError,
	LevelInfo,
}

// Logger is the main interface of this package, it's the only interface
// available to interact with every logger implementation, that's to keep the
// usage easily integrated in code and tests.
type Logger interface {
	// Tag returns a new Logger with the extra tag added to it. Tags are logged
	// on every line as a prefix of the message in the form [name:value]
	Tag(name string, value interface{}) Logger

	// Errorf will log a message applying args to the format string in the same
	// way fmt.Sprintf works but applying the [ERROR] prefix to the message
	Errorf(format string, args ...interface{})

	// Infof will log a message applying args to the format string in the same
	// way fmt.Sprintf works but applying the [INFO] prefix to the message
	Infof(format string, args ...interface{})
}

func LevelFromString(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "error":
		return LevelError, nil
	case "info":
		return LevelInfo, nil
	}

	var l Level
	return l, fmt.Errorf("not valid log level: %q", level)
}
