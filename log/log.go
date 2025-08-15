package log

import (
	"go.uber.org/zap"
)

var logger = zap.Must(zap.NewDevelopment()).Sugar()
var discard = zap.NewNop().Sugar()
var traceLogger = zap.Must(newTraceLogger()).Sugar()

var IsTrace = false

func GetLoggerWithSpan(span string) *zap.SugaredLogger {
	if span == "cpu" || span == "ppu" {
		if IsTrace {
			return traceLogger
		} else {
			return discard
		}
	} else {
		return logger
	}
}

func newTraceLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"famigom.trace",
	}
	return cfg.Build()
}
