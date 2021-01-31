package bdbstore

import (
	"io"

	"github.com/rs/zerolog"
)

func newLogger(logWriter io.Writer) logger {
	l := zerolog.New(logWriter).With().Str("component", "badgerbd").Timestamp().Logger()
	return logger{&l}
}

type logger struct {
	*zerolog.Logger
}

func (l logger) Errorf(f string, v ...interface{}) {
	l.Error().Msgf(f, v...)
}

func (l logger) Warningf(f string, v ...interface{}) {
	l.Warn().Msgf(f, v...)
}

func (l logger) Infof(f string, v ...interface{}) {
	l.Info().Msgf(f, v...)
}

func (l logger) Debugf(f string, v ...interface{}) {
	l.Debug().Msgf(f, v...)
}
