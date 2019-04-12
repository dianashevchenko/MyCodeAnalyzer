package metrics

import (
	"bufio"
	"os"
	"strings"
)

//NOP, LOC, HIT, NOM, CALL, NOC
//NProtM, NOPA, NOAV
type Metrics struct {
	packages          map[string]int
	methods           map[string]int
	variables         map[string]int
	variablesInMethod map[string]int
	protected         map[string]int
	public            map[string]int

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

const(
	equal = "=="
	assign = "="
	notEqual = "!="
)
var currentMethod = ""
var currentClass = ""
var bracketCount = 0
var previousLine = " "

func Count(files []string) Metrics {
	var dm Metrics
	dm.count(files)
	dm.Setting()
	//fmt.Printf("\nNop %v\nLOC %v\nNoc %v\nCall %v\nNoM %v\nHit %v\n"+
	//	"Noav %v\nNopa %v\nNprotm %v",
	//	len(dm.packages),
	//	dm.LOC,
	//	dm.NOC,
	//	dm.CALL,
	//	dm.NOM,
	//	float64(dm.HIT),
	//	dm.NOAV,
	//	dm.NOPA,
	//	dm.NProtM,
	//)
	return dm

}
func (dm *Metrics) Setting() {
	for _, v := range dm.methods {
		dm.CALL += v - 1
	}

	dm.NOP, dm.NOM = len(dm.packages), len(dm.methods)
	dm.HIT = dm.HIT / float64(dm.NOC)
	dm.NOAV, dm.NProtM = 0, 0

	delete(dm.variablesInMethod, "")
	delete(dm.variablesInMethod, "apply")
	delete(dm.variablesInMethod, "compare")

	for _, v := range dm.variablesInMethod {
		if dm.NOAV < v {
			dm.NOAV = v
		}
	}
	for _, v := range dm.protected {
		if dm.NProtM < v {
			dm.NProtM = v
		}
	}
	for _, v := range dm.public {
		if dm.NOPA < v {
			dm.NOPA = v
		}
	}
}

func (dm *Metrics) count(files []string) {
	dm.methods = make(map[string]int)
	dm.protected = make(map[string]int)
	dm.packages = make(map[string]int)
	dm.public = make(map[string]int)
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

			line = strings.TrimSpace(line)
			line = strings.Trim(line, "\t")
			if len(line) == 0 || dm.isComment(line, &isComment) {
				continue
			}
			if containsAnd(line, "{") && len(line) == 1 {
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

	if containsOr(line, `*/`, `//`) {
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
		vars := strings.Split(line, " ")
		for i, v := range vars {
			if v == "class" {
				currentClass = vars[i+1]
			}
		}
	}
	return true
}

func (dm *Metrics) checkForMethods(line string) bool {
	if containsOr(line, "import", "interface") {
		return false
	}

	//assume it is a method
	if containsAnd(line, "{", "(") && !containsKW(line) {
		bracketCount += strings.Count(line, "{")
		bracketCount -= strings.Count(line, "}")

		vars := strings.Split(line, " ")

		//is it a constructor?
		for i, v := range vars {
			if containsAnd(v, "(") {
				if i == 1 {
					dm.CALL++
					dm.methods[currentMethod]++
					return false
				}
				ind := strings.Index(v, "(")
				currentMethod = v[:ind]
			}
		}
		if containsOr(line, "private", "protected") {
			dm.protected[currentClass]++
		}
		dm.methods[currentMethod]++
		if !containsAnd(line, "()") {
			dm.variablesInMethod[currentMethod] += strings.Count(line, ",") + 1
		}
		return true
	}

	//assignment present => not a method
	if containsAnd(line, assign) && !containsAnd(line, equal, notEqual, ) {
		vars := strings.Split(line, " ")
		for i, v := range vars {
			if v == assign {
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
		if containsOr(line, "public") && containsAnd(line, "") {
			dm.public[currentClass]++
		}
		dm.methods[currentMethod]++
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

func containsKW(haystack string) bool {
	var kws = []string{"throws",
		"=", "new", "this", "while", "if", "switch", "for"}
	for _, v := range kws {
		if strings.Contains(haystack, v) {
			return true
		}
	}
	return false
}

