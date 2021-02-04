package settings

import (
	"io"
	slog "log"
	"strings"

	"github.com/mailchain/mailchain/cmd/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/internal/settings/values"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func logger(s values.Store) *Logger {
	l := &Logger{
		Path:       values.NewDefaultString(defaults.LogsPath(), s, "logger.path"),
		Level:      values.NewDefaultString("warn", s, "logger.level"),
		Console:    values.NewDefaultBool(true, s, "logger.console"),
		MaxBackups: values.NewDefaultInt(3, s, "logger.max-backups"),
		MaxAgeDays: values.NewDefaultInt(28, s, "logger.max-age-days"),
		MaxSizeMB:  values.NewDefaultInt(1, s, "logger.max-size-mb"),
		StackTrace: values.NewDefaultBool(false, s, "logger.stack-trace"),
	}

	return l
}

// Logger configuration element.
type Logger struct {
	Path       values.String
	Level      values.String
	Console    values.Bool
	MaxBackups values.Int
	MaxAgeDays values.Int
	MaxSizeMB  values.Int
	StackTrace values.Bool
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (l *Logger) Output() output.Element {
	return output.Element{
		FullName: "logger",
		Attributes: []output.Attribute{
			l.Path.Attribute(),
			l.Level.Attribute(),
			l.Console.Attribute(),
			l.MaxAgeDays.Attribute(),
			l.MaxBackups.Attribute(),
			l.MaxSizeMB.Attribute(),
			l.StackTrace.Attribute(),
		},
	}
}

func (l *Logger) Writer() io.Writer {
	return &lumberjack.Logger{
		Filename:   strings.Join([]string{l.Path.Get(), "mailchain.log"}, "/"),
		MaxSize:    l.MaxSizeMB.Get(),
		MaxBackups: l.MaxBackups.Get(),
		MaxAge:     l.MaxAgeDays.Get(),
	}
}

func (l *Logger) Init() {
	log.Logger = l.Produce()

	slog.SetOutput(&lumberjack.Logger{
		Filename:   strings.Join([]string{l.Path.Get(), "general.log"}, "/"),
		MaxSize:    l.MaxSizeMB.Get(),
		MaxBackups: l.MaxBackups.Get(),
		MaxAge:     l.MaxAgeDays.Get(),
	})
}

// Produce a logger.
func (l *Logger) Produce() zerolog.Logger {
	logger := zerolog.New(l.Writer())
	if l.Console.Get() {
		logger = zerolog.New(zerolog.MultiLevelWriter(l.Writer(), zerolog.NewConsoleWriter()))
	}

	if l.StackTrace.Get() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger = logger.With().Stack().Logger()
	}

	r := logger.Level(l.getLevel()).With().Caller().Timestamp().Logger()

	return r
}

func (l *Logger) getLevel() zerolog.Level {
	switch strings.ToLower(l.Level.Get()) {
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	}

	return zerolog.WarnLevel
}
