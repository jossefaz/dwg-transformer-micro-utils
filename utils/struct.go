package utils

type PickFile struct {
	Path string
	Result map[string]int
	From string
}


func SetResult(pFile *PickFile, resultKeys []string, resultVal []int, from string, path string){
	for i, k := range resultKeys {
		pFile.Result[k] = resultVal[i]
	}
	pFile.From = from
	pFile.Path = path
}

