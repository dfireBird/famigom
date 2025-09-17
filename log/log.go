package log

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap"
)

var logger = zap.Must(zap.NewDevelopment()).Sugar()
var IsTrace = false

func Logger() *zap.SugaredLogger {
	return logger
}

var w *bufio.Writer

func TraceLog(template string, args ...any) {
	if IsTrace {
		fmt.Fprintf(w, template, args...)
	}
}

func InitLogger(verbose bool) error {
	IsTrace = verbose

	f, err := os.Create("famigom.trace")
	if err != nil {
		return err
	}
	w = bufio.NewWriterSize(f, 512*1024*1024)
	return nil
}

func FlushLoggers() {
	w.Flush()
}
