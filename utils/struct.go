package utils

import "encoding/json"

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


func SetResultMessage(pFile *PickFile, resultKeys []string, resultVal []int, from string, path string) []byte{
	for i, k := range resultKeys {
		pFile.Result[k] = resultVal[i]
	}
	pFile.From = from
	pFile.Path = path
	res, err := json.Marshal(pFile)
	HandleError(err, "Cannot set result")
	return res
}
