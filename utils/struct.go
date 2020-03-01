package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type PickFile struct {
	Path string
	Result map[string]int
	From string
}

type Result struct {
	Success string
	Fail    string
	From    string
}

type Logger struct {
	Log *log.Logger
}


func SetResultMessage(pFile *PickFile, resultKeys []string, resultVal []int, from string, path string) ([]byte, error){
	for i, k := range resultKeys {
		pFile.Result[k] = resultVal[i]
	}
	pFile.From = from
	pFile.Path = path
	res, err := json.Marshal(pFile)
	return res, err
}
