package implement

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func LoggerTestSuite(t *testing.T, newLogger func(...LoggerOption) Logger) {
	for _, level := range allLevels {
		t.Run("Infof", func(t *testing.T) {
			var output bytes.Buffer
			testInfof(t, newLogger(WithLevel(level), Output(&output)), &output, level)
		})

		t.Run("Errorf", func(t *testing.T) {
			var output bytes.Buffer
			testErrorf(t, newLogger(WithLevel(level), Output(&output)), &output)
		})
	}

	t.Run("Tag", func(t *testing.T) {
		var output bytes.Buffer
		testTag(t, newLogger(WithLevel(LevelInfo), Output(&output)), &output)
	})

	t.Run("tags in lexicographic order", func(t *testing.T) {
		var output bytes.Buffer
		testTagsInLexicographicOrder(t, newLogger(WithLevel(LevelInfo), Output(&output)), &output)
	})
}

func testInfof(t *testing.T, logger Logger, output *bytes.Buffer, level Level) {
	format := "random message %d"
	arg := 1234
	message := fmt.Sprintf(format, arg)

	logger.Infof(format, arg)

	line := output.String()

	switch {
	case LevelInfo > level:
		if line != "" {
			t.Errorf("logger.Infof(%q, %v) = %q; want \"\"", format, arg, line)
		}

	default:
		if !strings.HasSuffix(line, message+"\n") {
			t.Errorf("logger.Infof(%q, %v) = %q; want %q", format, arg, line, message)
		}

		if !strings.Contains(line, "INFO") {
			t.Errorf("logger.Infof(%q, %v) = doesn't contain INFO level", format, arg)
		}
	}
}

func testErrorf(t *testing.T, logger Logger, output *bytes.Buffer) {
	format := "random message %d"
	arg := 1234
	message := fmt.Sprintf(format, arg)

	logger.Errorf(format, arg)

	line := output.String()

	if !strings.HasSuffix(line, message+"\n") {
		t.Errorf("logger.Errorf(%q, %v) = %q; want %q", format, arg, line, message)
	}

	if !strings.Contains(line, "ERROR") {
		t.Errorf("logger.Errorf(%q, %v) = doesn't contain ERROR level", format, arg)
	}
}

func testTag(t *testing.T, logger Logger, output *bytes.Buffer) {
	tags := map[string]interface{}{
		"aTagKey":           "string value",
		"a tag with spaces": "other string",
		"number tag":        12345,
	}
	for name, value := range tags {
		logger = logger.Tag(name, value)
	}

	format := "random message %d"
	arg := 1234
	message := fmt.Sprintf(format, arg)

	logger.Infof(format, arg)

	line := output.String()

	if !strings.HasSuffix(line, message+"\n") {
		t.Errorf("logger.Infof(%q, %v) = %q; want %q", format, arg, line, message)
	}
	if !strings.Contains(line, "INFO") {
		t.Errorf("logger.Infof(%q, %v) = doesn't contain INFO level", format, arg)
	}
	for name, value := range tags {
		want := fmt.Sprintf(`%s:%v`, name, value)
		if !strings.Contains(line, want) {
			t.Errorf("logger.Infof(%q, %v) = %q; want to contain tag %s:%v", format, arg, line, name, value)
		}
	}
}

func testTagsInLexicographicOrder(t *testing.T, logger Logger, output *bytes.Buffer) {
	tags := map[string]interface{}{
		"a": "string value",
		"A": 98765,
		"c": 12345,
		"B": "-",
		"b": "other string",
		"C": "!*%)!#",
	}
	for name, value := range tags {
		logger = logger.Tag(name, value)
	}

	order := make([]string, len(tags))
	for k, v := range tags {
		order = append(order, fmt.Sprintf(`[%s:%v] `, k, v))
	}
	sort.Strings(order)

	format := "random message %d"
	arg := 1234

	logger.Infof(format, arg)

	line := output.String()

	var index = 0
	for _, tag := range order {
		curr := strings.Index(line, tag)
		if curr < index {
			t.Errorf("logger.Infof(%q, %v) = %q; want tag %s in lexicographic order", format, arg, line, tag)
		}
		index = curr
	}
}
