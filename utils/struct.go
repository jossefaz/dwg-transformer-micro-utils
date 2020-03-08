package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type PickFile struct {
	Path string
	Result map[string]int
	From string
	To string
}

type Crud struct {
	Create string
	Retrieve string
	Update string
	Delete string
}

type DbQuery struct {
	dbType string
	Schema string
	Table string
	CrudT Crud
	Id []int
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


func SetResultMessage(pFile *PickFile, resultKeys []string, resultVal []int, from string, to string, path string) ([]byte, error){
	for i, k := range resultKeys {
		pFile.Result[k] = resultVal[i]
	}
	pFile.From = from
	pFile.Path = path
	pFile.To = to
	res, err := json.Marshal(pFile)
	return res, err
}
