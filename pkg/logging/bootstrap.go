package logging

import (
	"github.com/alxzoomer/clickhouse-explorer/pkg/appinfo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func Init(consoleLogger bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var wr io.Writer = os.Stderr
	if consoleLogger {
		wr = zerolog.ConsoleWriter{Out: os.Stderr}
	}
	log.Logger = log.
		With().
		Caller().
		Str("service", appinfo.Service).
		Str("version", appinfo.Version).
		Str("branch", appinfo.Branch).
		Str("buildtime", appinfo.BuildTime).
		Str("host", appinfo.Hostname).
		Logger().
		Output(wr)
}
