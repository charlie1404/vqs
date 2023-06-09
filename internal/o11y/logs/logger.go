package logs

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger() zerolog.Logger {
	// zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Logger = zerolog.New(os.Stdout).
	Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		// Int("pid", os.Getpid()).
		Logger()

	return Logger
}
