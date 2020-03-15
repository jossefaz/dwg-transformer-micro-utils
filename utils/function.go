package utils

import (
	"errors"
	"fmt"
	"os"
)

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

func GetEnv(envKey string) (string, error) {
	val, ok := os.LookupEnv(envKey)
	if !ok {
		return "", errors.New(fmt.Sprintf("%s not set\n", envKey))
	} else {
		return val, nil
	}

}