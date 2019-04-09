package globbing

import (
	"os"
	"path/filepath"
	"fmt"
)

var ganttproject = "E:/Univ/Sophomore/2s/EMPI/Chosen Projects/ganttproject-master/ganttproject"
var test = "E:/Univ/Sophomore/2s/EMPI/myAnalysis/New folder/"
var jasmine = "E:/Univ/Sophomore/2s/EMPI/Chosen Projects/Jasmine-0.5"
func GetFiles() ([]string, error){
	count := 0
	files, err := filepath.Glob("/*.java")
	if err != nil{
		panic(err)
	}
	_ = filepath.Walk(ganttproject,  func(path string, info os.FileInfo, err error) error {
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
	//fmt.Print(files)
	fmt.Printf("LENGTH IS %v", len(files))
	fmt.Printf("COUNT IS %v", count)
	return files, nil
}

func isJava(path string) bool{
	ext := filepath.Ext(path)
	if ext == ".java"{
		return true
	}
	return false
}