package utils

func HandleError(err error, msg string, logger *Logger) {
	if err != nil {
		logger.Log.Error("%s: %s", msg, err)
	}
}
func Itob(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}