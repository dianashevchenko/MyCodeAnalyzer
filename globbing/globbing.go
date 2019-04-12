package globbing

import (
	"os"
	"path/filepath"
)


func GetFiles(project string ) ([]string, error){
	files, err := filepath.Glob("/*.java")
	if err != nil{
		panic(err)
	}
	_ = filepath.Walk(project,  func(path string, info os.FileInfo, err error) error {
		if isJava(path){
			files  = append(files, path)
		}
		return nil
	})
	return files, nil
}

func isJava(path string) bool{
	if filepath.Ext(path) == ".java"{
		return true
	}
	return false
}