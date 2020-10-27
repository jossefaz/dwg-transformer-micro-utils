package logs

import (
	log "github.com/sirupsen/logrus"
	"github.com/yossefaz/dwg-transformer-micro-utils/utils"
	"io"
	"os"
)

func InitLogs(logFile string, level string) (utils.Logger, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return utils.Logger{}, err
	}
	l := log.New()
	l.SetOutput(io.MultiWriter(os.Stderr, file))
	l.SetFormatter(&log.JSONFormatter{})
	setLogLevel(level, l)
	return utils.Logger{
		Log: l,
	}, nil
}
func setLogLevel(level string, l *log.Logger) {

	switch level {
	case "DEBUG":
		l.SetLevel(log.DebugLevel)
	case "INFO":
		l.SetLevel(log.InfoLevel)
	case "WARN":
		l.SetLevel(log.WarnLevel)
	case "ERROR":
		l.SetLevel(log.FatalLevel)
	case "PANIC":
		l.SetLevel(log.PanicLevel)
	default:
		l.SetLevel(log.DebugLevel)
	}

}
