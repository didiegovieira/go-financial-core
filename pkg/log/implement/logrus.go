package implement

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	libruntime "github.com/didiegovieira/go-payments-core/pkg/runtime"
)

type loggerConfig struct {
	level  Level
	output io.Writer
}

var logLevelToLogrusLevel = map[Level]logrus.Level{
	LevelInfo:  logrus.InfoLevel,
	LevelError: logrus.ErrorLevel,
}

// NewLogrus creates a new Logger that uses https://github.com/sirupsen/logrus
// implementation. It configures that instance to defaults used by logistics
// applications.
func NewLogrus(options ...LoggerOption) Logger {
	var l loggerConfig
	defaultOptions := []LoggerOption{
		WithLevel(LevelInfo),
		Output(os.Stdout),
	}
	options = append(defaultOptions, options...)

	for _, option := range options {
		l = option(l)
	}

	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetOutput(l.output)
	logger.SetFormatter(logrusFormatter{
		timeFormat: "2006-01-02T15:04:05.1483386-00:00",
	})

	level, ok := logLevelToLogrusLevel[l.level]
	if !ok {
		level = logLevelToLogrusLevel[LevelInfo]
	}
	logger.SetLevel(level)

	return logrusAdapter{
		entry: logrus.NewEntry(logger),
	}
}

type logrusFormatter struct {
	timeFormat string
}

func (f logrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if entry == nil {
		return []byte{}, errors.New("nil logrus.Entry to format")
	}

	caller, file := prettifyCaller(entry.Caller)

	var buf bytes.Buffer
	_, err := fmt.Fprintf(&buf, `[%s] [%s] [%s] [%s] `,
		entry.Time.Format(f.timeFormat),
		strings.ToUpper(entry.Level.String()),
		caller, file,
	)
	if err != nil {
		return nil, err
	}

	tags := make([]string, len(entry.Data))
	for k, v := range entry.Data {
		tags = append(tags, fmt.Sprintf(`[%s:%v] `, k, v))
	}
	sort.Strings(tags)
	for _, tag := range tags {
		_, err := fmt.Fprint(&buf, tag)
		if err != nil {
			return nil, err
		}
	}

	n, err := buf.WriteString(entry.Message)
	if err != nil {
		return nil, err
	}
	if len(entry.Message) != n {
		return nil, fmt.Errorf("can't write the complete log message, want to write %d bytes, but wrote %d bytes", len(entry.Message), n)
	}

	n, err = buf.WriteRune('\n')
	if err != nil {
		return nil, err
	}
	if n != 1 {
		return nil, fmt.Errorf("can't write end of line, want to write 1 byte, but wrote %d bytes", n)
	}

	return buf.Bytes(), nil
}

func prettifyCaller(frame *runtime.Frame) (string, string) {
	// we need the next frame because logrus reports the caller as being
	// always logrusAdapter on this file
	frame = libruntime.Next(frame)
	if frame == nil {
		return "caller_not_found", "file_not_found"
	}
	fn := strings.Split(fmt.Sprintf("%s()", frame.Function), "/")
	return fn[len(fn)-1], fmt.Sprintf("%s:%d", path.Base(frame.File), frame.Line)
}

var _ Logger = logrusAdapter{}

type logrusAdapter struct {
	entry *logrus.Entry `di:"norecurse"`
}

func (l logrusAdapter) Tag(name string, value interface{}) Logger {
	l.entry = l.entry.WithField(name, value)
	return l
}

func (l logrusAdapter) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l logrusAdapter) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}
