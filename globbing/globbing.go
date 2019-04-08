package globbing

import (
	"os"
	"path/filepath"
	"fmt"
)

var dir = "E:/Univ/Sophomore/2s/EMPI/Chosen Projects/ganttproject-master/"
var dir2 = "E:/Univ/Sophomore/2s/EMPI/testing/"

func GetFiles() ([]string, error){
	files, err := filepath.Glob("/*.java")
	if err != nil{
		panic(err)
	}
	_ = filepath.Walk(dir2,  func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if isJava(path){
			fmt.Printf("visited java file: %q\n", path)
			files  = append(files, path)
		}
		return nil
	})
	fmt.Print(files)
	fmt.Printf("LENGTH IS %v", len(files))
	return files, nil
}

func isJava(path string) bool{
	ext := filepath.Ext(path)
	if ext == ".java"{
		return true
	}
	return false
}