package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func OpenFile(fname string) *os.File {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func ListFilesInDir(root string) []string {

	files := []string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		fi, err := os.Stat(path)

		if err != nil {
			fmt.Println(err)
			return nil
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			// do directory stuff
			fmt.Println("directory")
		case mode.IsRegular():
			// do file stuff
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
	return files
}
