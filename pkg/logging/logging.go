package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	LevelDebug Level = Level(zerolog.DebugLevel)
	LevelInfo  Level = Level(zerolog.InfoLevel)
	LevelWarn  Level = Level(zerolog.WarnLevel)
	LevelError Level = Level(zerolog.ErrorLevel)
)

type Level zerolog.Level

func New(
	level Level,
) *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(os.Stdout).
		Level(zerolog.Level(level)).
		With().Timestamp().Caller().Stack().Logger()

	return &logger
}
