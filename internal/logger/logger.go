package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Init(level string) error {
	var err error
	cfg := zap.NewProductionConfig()
	if cfg.Level, err = zap.ParseAtomicLevel(level); err != nil {
		return errors.Wrap(err, "parse atomic level")
	}

	if Log, err = cfg.Build(); err != nil {
		return errors.Wrap(err, "build")
	}

	return nil
}
