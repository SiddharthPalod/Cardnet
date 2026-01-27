package config

import "sync/atomic"

type RuntimeConfig struct {
	VelocityHighTx1Min int
	VelocityHighScore  int
}

var activeConfig atomic.Value

func Load(cfg RuntimeConfig) {
	activeConfig.Store(cfg)
}

func Get() RuntimeConfig {
	return activeConfig.Load().(RuntimeConfig)
}
