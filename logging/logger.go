package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
)

type (
	HandleType  uint8
	Environment uint8
)

const (
	ConsoleHandler HandleType = iota
	TextHandler
	JSONHandler
)

const (
	DEVELOPMENT Environment = iota
	PRODUCTION
	RELEASE
)

const (
	_defaultSentryFlushTimeout = 1 * time.Second
)

type Log struct {
	skipCaller int
	slog       *slog.Logger
	sentry     *sentry.Client
}

type SentryConfig struct {
	DSN                string
	ServerName         string
	Environment        Environment
	AttachStacktrace   bool
	SampleRate         float64
	EnableTracing      bool
	TracesSampleRate   float64
	ProfilesSampleRate float64
	MaxBreadcrumbs     int
	Release            string
	Dist               string
	HTTPProxy          string
	HTTPSProxy         string
	MaxErrorDepth      int
	Debug              bool
}

type Options struct {
	Development  bool // Development add development details of machine
	Debug        bool // Debug show debug devel message
	EnableCaller bool // EnableCaller show caller in line code
	SkipCaller   int  // SkipCaller skip caller level of CallerFrames
	// https://github.com/golang/go/issues/59145#issuecomment-1481920720
	Sentry *SentryConfig // Sentry enable sentry with specific configuration
}

type Logger interface {
	Debug(toSentry bool, msg string, args ...any)
	DebugContext(ctx context.Context, toSentry bool, msg string, args ...any)
	Info(toSentry bool, msg string, args ...any)
	InfoContext(ctx context.Context, toSentry bool, msg string, args ...any)
	Warn(toSentry bool, msg string, args ...any)
	WarnContext(ctx context.Context, toSentry bool, msg string, args ...any)
	Error(toSentry bool, msg string, args ...any)
	ErrorContext(ctx context.Context, toSentry bool, msg string, args ...any)
	Fatal(toSentry bool, msg string, args ...any)
	FatalContext(ctx context.Context, toSentry bool, msg string, args ...any)
	Log(ctx context.Context, toSentry bool, level slog.Level, msg string, args ...any)

	GetSentryClient() *sentry.Client
}

func New(
	handler HandleType,
	loggerOption Options,
) (Logger, error) {
	log := new(Log)
	logger := slog.Default()
	slogHandlerOpt := new(slog.HandlerOptions)
	slogHandlerOpt.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		return a
	}

	if loggerOption.Debug {
		slogHandlerOpt.Level = slog.LevelDebug
	}

	if loggerOption.EnableCaller {
		slogHandlerOpt.AddSource = true
	}

	if loggerOption.Sentry != nil {
		client, err := sentry.NewClient(setSentryOptions(loggerOption.Sentry))
		if err != nil {
			return nil, err
		}

		log.sentry = client
	}

	switch handler {
	case JSONHandler:
		logger = slog.New(slog.NewJSONHandler(os.Stderr, slogHandlerOpt))
	case TextHandler:
		logger = slog.New(slog.NewTextHandler(os.Stderr, slogHandlerOpt))
	case ConsoleHandler:
		logger = slog.New(NewConsoleHandler(slogHandlerOpt))
	}

	if loggerOption.Development {
		buildInfo, _ := debug.ReadBuildInfo()
		logger = logger.With(slog.Group("debug_info",
			slog.String("go_version", buildInfo.GoVersion),
			slog.Int("pid", os.Getpid()),
			slog.String("os", runtime.GOOS),
			slog.String("os_arch", runtime.GOARCH),
		))
	}

	log.slog = logger
	log.skipCaller = loggerOption.SkipCaller

	return log, nil
}

func (l *Log) Debug(toSentry bool, msg string, keyValues ...any) {
	l.Log(context.Background(), toSentry, slog.LevelDebug, msg, keyValues...)
}

func (l *Log) DebugContext(ctx context.Context, toSentry bool, msg string, keyValues ...any) {
	l.Log(ctx, toSentry, slog.LevelDebug, msg, keyValues...)
}

func (l *Log) Info(toSentry bool, msg string, keyValues ...any) {
	l.Log(context.Background(), toSentry, slog.LevelInfo, msg, keyValues...)
}

func (l *Log) InfoContext(ctx context.Context, toSentry bool, msg string, keyValues ...any) {
	l.Log(ctx, toSentry, slog.LevelInfo, msg, keyValues...)
}

func (l *Log) Warn(toSentry bool, msg string, keyValues ...any) {
	l.Log(context.Background(), toSentry, slog.LevelWarn, msg, keyValues...)
}

func (l *Log) WarnContext(ctx context.Context, toSentry bool, msg string, keyValues ...any) {
	l.Log(ctx, toSentry, slog.LevelWarn, msg, keyValues...)
}

func (l *Log) Error(toSentry bool, msg string, keyValues ...any) {
	l.Log(context.Background(), toSentry, slog.LevelError, msg, keyValues...)
}

func (l *Log) ErrorContext(ctx context.Context, toSentry bool, msg string, keyValues ...any) {
	l.Log(ctx, toSentry, slog.LevelError, msg, keyValues...)
}

func (l *Log) Fatal(toSentry bool, msg string, keyValues ...any) {
	defer os.Exit(1)
	l.Log(context.Background(), toSentry, slog.LevelError, msg, keyValues...)
}

func (l *Log) FatalContext(ctx context.Context, toSentry bool, msg string, keyValues ...any) {
	defer os.Exit(1)
	l.Log(ctx, toSentry, slog.LevelError, msg, keyValues...)
}

func (l *Log) Log(ctx context.Context, toSentry bool, level slog.Level, msg string, keyValues ...any) {
	var pcs [1]uintptr
	runtime.Callers(l.skipCaller, pcs[:])
	rec := slog.NewRecord(time.Now(), level, msg, pcs[0])
	rec.Add(keyValues...)

	if toSentry && l.sentry != nil {
		defer l.sentry.Flush(_defaultSentryFlushTimeout)
		sentryLevel := sentry.LevelInfo
		switch level {
		case slog.LevelWarn:
			sentryLevel = sentry.LevelWarning
		case slog.LevelError:
			sentryLevel = sentry.LevelError
		case slog.LevelDebug:
			sentryLevel = sentry.LevelDebug
		}
		event := l.sentry.EventFromMessage(fmt.Sprint(msg, keyValues), sentryLevel)
		l.sentry.CaptureEvent(event, nil, nil)
		rec.Add(
			"sent_to_sentry", true,
			"sentry_event_id", event.EventID,
		)
	}

	_ = l.slog.Handler().Handle(ctx, rec)
}

func (l *Log) GetSentryClient() *sentry.Client {
	return l.sentry
}

func (e Environment) String() string {
	switch e {
	case DEVELOPMENT:
		return "development"
	case PRODUCTION:
		return "production"
	case RELEASE:
		return "release"
	}

	return ""
}

func setSentryOptions(cfg *SentryConfig) sentry.ClientOptions {
	opt := sentry.ClientOptions{
		Dsn:              cfg.DSN,
		Debug:            cfg.Debug,
		AttachStacktrace: cfg.AttachStacktrace,
		EnableTracing:    cfg.EnableTracing,
	}

	if len(cfg.Dist) != 0 {
		opt.Dist = cfg.Dist
	}

	if len(cfg.ServerName) != 0 {
		opt.ServerName = cfg.ServerName
	}

	if len(cfg.HTTPProxy) != 0 {
		opt.HTTPProxy = cfg.HTTPProxy
	}

	if len(cfg.HTTPSProxy) != 0 {
		opt.HTTPSProxy = cfg.HTTPSProxy
	}

	if len(cfg.Release) != 0 {
		opt.Release = cfg.Release
	}

	if cfg.MaxErrorDepth != 0 {
		opt.MaxErrorDepth = cfg.MaxErrorDepth
	}

	if cfg.TracesSampleRate != 0 {
		opt.TracesSampleRate = cfg.TracesSampleRate
	}

	if cfg.EnableTracing {
		opt.TracesSampleRate = 1.0
		if cfg.TracesSampleRate != 0 {
			opt.TracesSampleRate = cfg.TracesSampleRate
		}
	}

	if cfg.ProfilesSampleRate != 0 {
		opt.ProfilesSampleRate = cfg.ProfilesSampleRate
	}

	if cfg.SampleRate != 0 {
		opt.SampleRate = cfg.SampleRate
	}

	return opt
}
