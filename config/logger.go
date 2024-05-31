package config

import "github.com/phuslu/log"

func InitLogger() {
	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     1,
		TimeField:  "date",
		TimeFormat: "2006-01-02",
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}
}
