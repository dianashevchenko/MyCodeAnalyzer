package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	//"time"

	//"time"
)

//NOP, LOC, HIT, NOM, CALL, NOC
//NProtM, NOPA, NOAV
type Metrics struct {
	packages          map[string]int
	methods           map[string]int
	variables         map[string]int
	variablesInMethod map[string]int

	NOP    int
	LOC    int
	NOC    int
	CALL   int
	NOM    int
	HIT    float64
	NOAV   int
	NOPA   int
	NProtM int
}

//var notAllowed = []string{
//	"abstract", "assert", "boolean", "break",
//	"byte", "case", "catch", "char",
//	"class", "const", "continue", "default",
//	"do", "double", "else", "enum",
//	"extends", "final", "finally", "float",
//	"for", "goto", "if", "implements",
//	"import", "instanceof", "int", "interface",
//	"long", "native", "new", "package",
//	"private", "protected", "public", "return",
//	"short", "static", "strictfp", "super",
//	"switch", "synchronized", "this", "throw",
//	"throws", "transient", "try", "void",
//	"volatile", "while", "true", "false",
//	"null",
//}

var currentMethod = ""
var bracketCount = 0
var totalVariables = 0
var previousLine = " "

func Count(files []string) Metrics {
	var dm Metrics
	dm.count(files)
	for _, v := range dm.methods {
		dm.CALL += v - 1
	}
	//fmt.Printf("HIT + %v", (dm.HIT))
	dm.NOP = len(dm.packages)
	dm.NOM = len(dm.methods)
	dm.HIT = dm.HIT / float64(dm.NOC)
	dm.NOAV = 0
	totalVariables = dm.variablesInMethod[""]
	delete(dm.variablesInMethod, "")
	delete(dm.variablesInMethod, "apply")
	delete(dm.variablesInMethod, "compare")
	for _, v := range dm.variablesInMethod {

		//time.Sleep(time.Second)
		//fmt.Println("Noav + "+ string(v))
		if dm.NOAV < v {
			//fmt.Printf("%v\n", k)
			//fmt.Printf("NOAV = %v  v = %v\n", dm.NOAV, v)
			dm.NOAV = v
		}
		//dm.NOAV += v - 1
	}

	fmt.Printf("\nNop %v\nLOC %v\nNoc %v\nCall %v\nNoM %v\nHit %v\n"+
		"Noav %v\nNopa %v\nNprotm %v",
		len(dm.packages),
		dm.LOC,
		dm.NOC,
		dm.CALL,
		dm.NOM,
		float64(dm.HIT),
		dm.NOAV,
		dm.NOPA,
		dm.NProtM,
	)
	return dm

}
func (dm *Metrics) count(files []string) {
	dm.methods = make(map[string]int)
	dm.packages = make(map[string]int)
	dm.variables = make(map[string]int)
	dm.variablesInMethod = make(map[string]int)

	for i := range files {
		file, err := os.Open(files[i])
		if err != nil {
			panic(err)
		}

		isComment := false
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			//if containsAnd(line, "compare"){
			//	fmt.Println(files[i])
			//	time.Sleep(time.Second*2)
			//}

			line = strings.TrimSpace(line)
			line = strings.Trim(line, "\t")
			if len(line) ==0 || dm.isComment(line, &isComment) {
				continue
			}
			if containsAnd(line, "{") && len(line) ==1{
				line = previousLine + "{"
			}
			dm.LOC++
			dm.runChecks(line)
			previousLine = line
			_ = file.Close()
		}
	}

}
func (dm *Metrics) isComment(line string, isComment *bool) bool {

	if strings.Contains(line, `//`) {
		return true
	}
	if strings.Contains(line, `*/`) {
		*isComment = false
		return true
	}
	if strings.Contains(line, `/*`) {
		*isComment = true
		return true
	}
	if *isComment == true {
		return true
	}
	return false
}
func (dm *Metrics) runChecks(line string) {

	dm.checkForClasses(line)

	dm.checkForPackages(line)

	dm.checkForMethods(line)

}

func (dm *Metrics) checkForClasses(line string) bool {
	if strings.Contains(line, "extends") {
		dm.HIT++
	}
	if strings.Contains(line, "class") && !strings.Contains(line, ".class") {
		dm.NOC++

	}
	return true
}

func (dm *Metrics) checkForMethods(line string) bool {
	if containsOr(line, "import", "interface") {
		return false
	}



	//assume it is a method
	if containsAnd(line, "{", "(") && !containsOr(line, "throws", "=", "new", "this", "while", "if",
		"switch", "for") {
		bracketCount += strings.Count(line, "{")
		bracketCount -= strings.Count(line, "}")


		vars := strings.Split(line, " ")

		//is it a constructor?

		for i, v := range vars {
			if containsAnd(v, "(") {
				if i == 1 {
					dm.CALL++
					return false

				}
				ind := strings.Index(v, "(")
				currentMethod = v[:ind]

			}

		}
		if containsOr(line, "private", "protected") {
			dm.NProtM++
		}
		dm.methods[currentMethod]++
		if containsAnd(line, "()") {
			//there are no input parameters

		} else {
			dm.variablesInMethod[currentMethod] += strings.Count(line, ",") + 1
		}
		return true
	}

	//there is some assignment => not a method
	if containsAnd(line, "=") && !containsAnd(line, "==", "!=", ) {
		vars := strings.Split(line, " ")
		for i, v := range vars {
			if v == "=" {
				var variable string
				if containsAnd(vars[i-1], ".") {
					variable = vars[i-1][:strings.Index(vars[i-1], ".")]
				}
				if dm.variables[variable] == 0 {
					dm.variables[variable]++
					dm.variablesInMethod[currentMethod]++
				}

			}
		}
		if containsOr(line, "public", "static") {
			dm.NOPA++
		}
		return false
	}

	if bracketCount == 0 {
		currentMethod = ""
	}

	return true
}

func (dm *Metrics) checkForPackages(line string) bool {

	if strings.Contains(line, "package") {
		dm.packages[line]++
		return true
	}
	return false
}

func containsAnd(haystack string, needle ...string) bool {
	for _, v := range needle {
		if !strings.Contains(haystack, v) {
			return false
		}
	}
	return true
}
func containsOr(haystack string, needle ...string) bool {
	for _, v := range needle {
		if strings.Contains(haystack, v) {
			return true
		}
	}
	return false
}

func sleep2() {
	time.Sleep(time.Second * 2)
}
