package implement

import "io"

type LoggerOption func(loggerConfig) loggerConfig

func WithLevel(level Level) LoggerOption {
	return func(l loggerConfig) loggerConfig {
		l.level = level
		return l
	}
}

func Output(w io.Writer) LoggerOption {
	return func(l loggerConfig) loggerConfig {
		l.output = w
		return l
	}
}
