package utils

func HandleError(err error, msg string, logger *Logger) {
	if err != nil {
		logger.Log.Error("%s: %s", msg, err)
	}
}
