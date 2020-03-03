package logs

import (
	"github.com/yossefazoulay/go_utils/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLogs(logFile string, level string) utils.Logger {
	var file, err = os.Create(logFile)
	if err != nil {
		panic(err)
	}
	l := log.New()
	l.SetOutput(file)
	l.SetFormatter(&log.JSONFormatter{})
	setLogLevel(level, l)
	return utils.Logger{
		Log:l,
	}
}
func setLogLevel(level string, l *log.Logger) {

	switch level {
	case "DEBUG" :
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
