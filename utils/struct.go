package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type PickFile struct {
	Path string
	Result map[string]int
}

type DbQuery struct {
	dbType string
	Schema string
	Table string
	CrudT string
	Id map[string]interface{}
	ORMKeyVal map[string]interface{}
}

type Result struct {
	Success string
	Fail    string
	From    string
}

type Logger struct {
	Log *log.Logger
}


func SetResultMessage(pFile *PickFile, resultKeys []string, resultVal []int, path string) ([]byte, error){
	for i, k := range resultKeys {
		pFile.Result[k] = resultVal[i]
	}
	pFile.Path = path
	res, err := json.Marshal(pFile)
	return res, err
}
