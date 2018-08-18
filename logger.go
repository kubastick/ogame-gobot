package main

import (
	"github.com/op/go-logging"
	"os"
	"runtime"
)

func setupLogger() *logging.Logger {
	var log = logging.MustGetLogger("Ogame Bot")
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{color:reset}%{message}`,
	)

	//Windows terminal does not support all utf-8 characters
	if runtime.GOOS == "windows" {
		format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} => %{level:.4s} %{color:reset}%{message}`)
	}

	backend2 := logging.NewLogBackend(os.Stdout, "", 0)
	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Set the backends to be used.
	logging.SetBackend(backend2Formatter)
	log.Info("Logger ready")
	return log
}
