package util

import (
	"github.com/arthurkiller/rollingwriter"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// Log 日志
var Log *zerolog.Logger

// GetLog 获取使用的日志
func GetLog() *zerolog.Logger {
	if Log == nil {
		Log = NewDefaultLog()
	}
	return Log
}

// SetLog 设置使用的日志
func SetLog(log *zerolog.Logger) {
	if log == nil {
		log = NewDefaultLog()
	}
	Log = log
}

// NewDefaultLog 创建默认日志
func NewDefaultLog() *zerolog.Logger {
	return newLogByCfg(nil)
}

func newLogByCfg(cfg *rollingwriter.Config) *zerolog.Logger {
	if cfg == nil {
		cfg = Config.NewDefaultLogCfg()
	}

	w, err := rollingwriter.NewWriterFromConfig(cfg)
	if err != nil {
		// 如果没创建日志成功，就使用zerolog的log
		zlog.Error().Err(err).Msg("使用zerolog的默认log")
		return &zlog.Logger
	}
	log := zerolog.New(w).With().Timestamp().Logger()
	return &log
}
