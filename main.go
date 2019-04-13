package main

import (
	"./globbing"
	"./metrics"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

func main() {

	address := os.Args[1:][0]
	project, err := url.QueryUnescape(address)

	files, err := globbing.GetFiles(project)
	if err != nil {
		fmt.Println(err)
	}

	result, str := metrics.Count(files)
	data, err := json.MarshalIndent(result, "", " ")
///metricsAnalyzer/result.json
	resultDir := project + "/metricsAnalyzer"
	_ = os.Mkdir(resultDir, 0755)
	f, err := os.Create(resultDir + "/result.json")
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
	_, err = f.Write(data)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	f.Close()
	f, err = os.Create(resultDir + "/result.txt")

	_, err = f.Write([]byte(str))
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	f.Close()

	//defer time.Sleep(10 * time.Second)
}

/*
func display(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("layout.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		panic(err)
	}
	//fmt.Fprintf(w, "fvfv %v", 0)
}
*/