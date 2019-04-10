package globbing

import (
	"os"
	"path/filepath"
	"fmt"
)


func GetFiles(project string ) ([]string, error){
	count := 0
	files, err := filepath.Glob("/*.java")
	if err != nil{
		panic(err)
	}
	_ = filepath.Walk(project,  func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if isJava(path){
			files  = append(files, path)
			count++
		}
		return nil
	})
	return files, nil
}

func isJava(path string) bool{
	ext := filepath.Ext(path)
	if ext == ".java"{
		return true
	}
	return false
}