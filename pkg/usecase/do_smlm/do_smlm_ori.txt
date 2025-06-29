package main

import (
	_ds "smlm_automation/pkg/usecase/do_smlm"
	_const "smlm_automation/pkg/util/constants"
	log "smlm_automation/pkg/util/logger"
)

func main() {
	loggerConfig := log.Config{
		Level:             _const.Loglevel,
		TimestampFormat:   _const.TimestampFormat,
		StdoutEnabled:     _const.StdoutEnabled,
		FilePath:          _const.LogFileCheckFabric,
		MaxSize:           _const.MaxSize,
		StacktraceEnabled: _const.StacktraceEnabled,
		ServiceName:       "checkfabric",
		EnableFileLogs:    _const.EnableFileLogs,
	}
	// Initialize logger object
	zapConfig, _ := log.NewConfig(&loggerConfig)
	logger, _, _ := log.New(&loggerConfig, zapConfig)
	logger.Debug("starting do_smlm main")
	doSmlm := _ds.NewDoSmlm(logger)
	doSmlm.StartDoSmlm()
	// Keep the program running
	select {}
}
