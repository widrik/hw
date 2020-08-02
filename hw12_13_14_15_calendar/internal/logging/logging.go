package logging

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ErrorLvl string = "error"
	WarnLvl  string = "warn"
	InfoLvl  string = "info"
	DebugLvl string = "debug"
)

func Init(lvl string, file string) error {
	var atom zap.AtomicLevel

	switch lvl {
	case ErrorLvl:
		atom = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case WarnLvl:
		atom = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case InfoLvl:
		atom = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case DebugLvl:
		atom = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	default:
		return errors.New("level does not exist")
	}

	if file == "" {
		return errors.New("file is wrong")
	}

	config := zap.NewProductionConfig()

	config.Level = atom
	config.Encoding = "json"
	config.OutputPaths = append(config.OutputPaths, file)

	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"

	Logger, err := config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(Logger)

	return nil
}
