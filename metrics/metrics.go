package metrics

import (
	"bufio"
	"os"
	"strings"
	//"time"

	//"time"
)

//NOP, LOC, HIT, NOM, CALL, NOC
//NProtM, NOPA, NOAV
type DirectMetrics struct {
	packages map[string]int
	methods  map[string]int

	NOP      int
	LOC      int
	NOC      int
	CALL     int
	NOM      int
	HIT      float64
	NOAV     int
	NOPA     int
	NProtM   int
}

func Count(files []string) DirectMetrics {
	var dm DirectMetrics
	dm.count(files)
	for _, v := range dm.methods {
		dm.CALL += v
	}
	dm.NOP = len(dm.packages)
	dm.NOM = len(dm.methods)
	dm.HIT = dm.HIT / float64(dm.NOC)
	/*fmt.Printf("\nNop %v\nLOC %v\nNoc %v\nCall %v\nNoM %v\nHit %v\n"+
		"Noav %v\nNopa %v\nNprotm %v",
		len(dm.packages),
		dm.LOC,
		dm.NOC,
		dm.CALL,
		len(dm.methods),
		float64(dm.HIT/float64(dm.NOC)),
		dm.NOAV,
		dm.NOPA,
		dm.NProtM,
	)*/
	return dm

}
func (dm *DirectMetrics) count(files []string) {
	dm.methods = make(map[string]int)
	dm.packages = make(map[string]int)
	for i := range files {
		file, err := os.Open(files[i])
		defer file.Close()

		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		//isComment := false
		for scanner.Scan() {
			line := scanner.Text()

			dm.LOC++
			//if dm.isComment(line, &isComment){
			//	//fmt.Println(line)
			//	continue
			//	//time.Sleep(time.Second)
			//}

			dm.checkForClasses(line)
			dm.checkForMethods(line)
			dm.checkForPackages(line)

		}
	}

}
func (dm *DirectMetrics) isComment(line string, isComment *bool) bool {

	if strings.Contains(line, `//`) {
		return true
	}
	if strings.Contains(line, `/*`) && strings.Contains(line, `*/`) {
		*isComment = false
		return false
	}
	if strings.Contains(line, `/*`) {
		*isComment = true
		return true
	}
	if strings.Contains(line, `*/`) {
		*isComment = false
		return true
	}
	if *isComment == true {

		return true
	}
	return false
}

/*func (dm *DirectMetrics) checkForClasses(line string) {
	if strings.Contains(line, "class") {
		arr := strings.Split(line, " ")
		for i, v := range arr {
			if v == "class" {

				if len(dm.NOC[arr[i+1]]) == 0 {

				}
				dm.NOC = make(map[string][]string)
				dm.NOC[arr[i+1]] = []string{}
				i++
			}
			if v == "extends" {
				dm.NOC[arr[i+1]] = append(dm.NOC[arr[i+1]], arr[i-1])
				i++
			}
		}
	}
}*/

func (dm *DirectMetrics) checkForClasses(line string) {
	if strings.Contains(line, "extends") {
		dm.HIT++
	}
	if strings.Contains(line, "class") && !strings.Contains(line, ".class") {
		dm.NOC++
	}
}
func (dm *DirectMetrics) checkForMethods(line string) {
	if strings.Contains(line, "class") ||
		strings.Contains(line, "extends") ||
		strings.Contains(line, "package") ||
		strings.Contains(line, "import") ||
		strings.Contains(line, "interface") {
		return
	}
	if strings.Contains(line, "private") || strings.Contains(line, "protected") {
		dm.NProtM++
		//fmt.Println(line)
		//time.Sleep(time.Second)
	}
	if strings.Contains(line, "public") || strings.Contains(line, "static") {
		dm.NOPA++
	}
	if strings.Contains(line, "=") ||
		!strings.Contains(line, "(") ||
		strings.Contains(line, ";") && !strings.Contains(line, "{") {
		return
	}

	if strings.Contains(line, "Override") {
		dm.methods["Override"]++
	}
	//time.Sleep(2 * time.Second)
	ind := strings.Index(line, "(")
	if ind != -1 {
		line = string(line[:ind])
	}
	line = strings.Trim(line, "\t .};")
	arr := strings.Split(line, " ")
	//fmt.Print(arr)

	dm.methods[arr[len(arr)-1]]++

}

func (dm *DirectMetrics) checkForPackages(line string) {

	if strings.Contains(line, "package") {
		dm.packages[line]++
		return
	}

}
