package logging

import (
	"context"

	"github.com/rs/zerolog"
)

func Ctx(ctx context.Context) *zerolog.Logger {
	logger := zerolog.Ctx(ctx)
	if logger != nil {

		return logger
	}

	return New(LevelDebug)
}
