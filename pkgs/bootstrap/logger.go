package bootstrap

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.aporeto.io/a3s/internal/conf"
	"go.aporeto.io/addedeffect/tracer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConfigureLogger configures the logging subsystem.
func ConfigureLogger(serviceName string, cfg conf.LoggingConf) tracer.CloseRecorderHandler {

	var err error

	configureZap(serviceName, cfg.LogLevel, cfg.LogFormat)

	f, err := tracer.ConfigureTracerWithURL(cfg.LogTracerURL, serviceName)
	if err != nil {
		zap.L().Warn("Unable to configure the OpenTracing", zap.Error(err))
	}

	if f != nil {
		zap.L().Info("OpenTracing enabled", zap.String("server", cfg.LogTracerURL))
	}

	return f
}

func configureZap(name string, level string, format string) {

	var config zap.Config

	switch format {
	case "json":
		config = getJSONConfig()
	case "stackdriver":
		config = getStackdriverConfig()
	default:
		config = getDefaultConfig()
	}

	config.Level = levelToZapLevel(level)

	enc, err := getEncoder(config)
	if err != nil {
		panic(err)
	}

	core := zapcore.NewCore(enc, zapcore.Lock(os.Stderr), config.Level)

	var logger *zap.Logger
	if name != "" {
		logger = zap.New(core, zap.Fields(zap.String("srv", name)))
	} else {
		logger = zap.New(core)
	}

	zap.ReplaceGlobals(logger)

	go handleElevationSignal(config)
}

func getJSONConfig() zap.Config {

	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.CallerKey = "c"
	config.EncoderConfig.LevelKey = "l"
	config.EncoderConfig.MessageKey = "m"
	config.EncoderConfig.NameKey = "n"
	config.EncoderConfig.TimeKey = "t"

	return config
}

func getStackdriverConfig() zap.Config {

	config := zap.NewProductionConfig()
	config.EncoderConfig.LevelKey = "severity"
	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		}
	}

	return config
}

func getDefaultConfig() zap.Config {

	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true
	config.DisableCaller = true
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return config
}

func levelToZapLevel(level string) zap.AtomicLevel {

	switch level {
	case "trace", "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func getEncoder(c zap.Config) (zapcore.Encoder, error) {

	switch c.Encoding {
	case "json":
		return zapcore.NewJSONEncoder(c.EncoderConfig), nil
	case "console":
		return zapcore.NewConsoleEncoder(c.EncoderConfig), nil
	default:
		return nil, errors.New("unknown encoding")
	}
}
func handleElevationSignal(cfg zap.Config) {

	defaultLevel := cfg.Level
	var elevated bool

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	for s := range c {
		if s == syscall.SIGINT {
			return
		}
		elevated = !elevated

		if elevated {
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			l, _ := cfg.Build()
			zap.ReplaceGlobals(l)
			zap.L().Info("Log level elevated to debug")
		} else {
			zap.L().Info("Log level restored to original configuration", zap.Stringer("level", defaultLevel.Level()))
			cfg.Level = defaultLevel
			l, _ := cfg.Build()
			zap.ReplaceGlobals(l)
		}
	}
}
