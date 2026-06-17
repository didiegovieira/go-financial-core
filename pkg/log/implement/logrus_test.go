package implement

import (
	"testing"
)

func Test_logrusAdapter(t *testing.T) {
	LoggerTestSuite(t, NewLogrus)
}
