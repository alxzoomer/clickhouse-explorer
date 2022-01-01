package logging

import (
	"github.com/rs/zerolog/log"
	stdlog "log"
)

type ZeroLogErrorLogAdapter struct{}

func (l *ZeroLogErrorLogAdapter) Write(p []byte) (n int, err error) {
	log.Error().Msg(string(p))
	return len(p), nil
}

func NewErrorLog() *stdlog.Logger {
	return stdlog.New(&ZeroLogErrorLogAdapter{}, "", 0)
}
