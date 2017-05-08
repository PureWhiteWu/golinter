package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

type lintFunc func(*os.File) (Result, error)

func lintJava(tmpFile *os.File) (Result, error) {
	cmd := exec.Command("java", "-jar", "linters/java/checkstyle-7.7-all.jar", "-c", "linters/java/sun_checks.xml", tmpFile.Name())
	lintOutput, _ := cmd.CombinedOutput()
	re, _ := regexp.Compile(`(\[ERROR].*?)\n`)
	tmp := re.FindAllStringSubmatch(string(lintOutput), -1)

	// replace unnecessary information with ""
	var lintErrors []string
	for _, strSlice := range tmp {
		re, _ = regexp.Compile("/.*?java:")
		t := re.ReplaceAllString(strSlice[1], "")
		lintErrors = append(lintErrors, t)
	}

	result := Result{
		ErrorNum: len(lintErrors),
		Errors:   lintErrors,
	}
	return result, nil
}

func lintCpp(tmpFile *os.File) (Result, error) {
	cmd := exec.Command("python", "linters/cpp/cpplint.py", "--filter=-legal/copyright", tmpFile.Name())
	lintOutput, _ := cmd.CombinedOutput()
	re, _ := regexp.Compile(`cpp:(.*?)\n`)
	tmp := re.FindAllStringSubmatch(string(lintOutput), -1)

	var lintErrors []string
	for _, strSlice := range tmp {
		lintErrors = append(lintErrors, strSlice[1])
	}

	result := Result{
		ErrorNum: len(lintErrors),
		Errors:   lintErrors,
	}
	return result, nil
}

func lintPython(tmpFile *os.File) (Result, error) {
	cmd := exec.Command("flake8", tmpFile.Name())
	lintOutput, _ := cmd.CombinedOutput()
	re, _ := regexp.Compile("/.*?py:(.*?)\n")
	tmp := re.FindAllStringSubmatch(string(lintOutput), -1)

	var lintErrors []string
	for _, strSlice := range tmp {
		lintErrors = append(lintErrors, strSlice[1])
	}

	result := Result{
		ErrorNum: len(lintErrors),
		Errors:   lintErrors,
	}
	return result, nil
}

func lintJavascript(tmpFile *os.File) (Result, error) {
	os.Chdir("./linters/javascript/")
	cmd := exec.Command("npm", "run", "lint", tmpFile.Name())
	lintOutput, _ := cmd.CombinedOutput()
	re, _ := regexp.Compile(`(\d+?:\d+?.*?)\n`)
	tmp := re.FindAllStringSubmatch(string(lintOutput), -1)

	var lintErrors []string
	for _, strSlice := range tmp {
		lintErrors = append(lintErrors, strSlice[1])
	}

	result := Result{
		ErrorNum: len(lintErrors),
		Errors:   lintErrors,
	}

	return result, nil
}

// dispatch linters
func dispatch(code Code) lintFunc {
	switch code.Language {
	case "java":
		return lintJava
	case "cpp":
		return lintCpp
	case "python":
		return lintPython
	case "javascript":
		return lintJavascript
	default:
		return nil
	}
}

// create the temp file
func createTempFile(code Code) (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "lint_")
	if err != nil {
		return nil, err
	}
	defer tmpfile.Close()

	_, err = tmpfile.WriteString(code.Source)
	if err != nil {
		return nil, err
	}

	err = tmpfile.Sync()
	if err != nil {
		return nil, err
	}
	suffix := "." + languageSuffix[code.Language]
	err = os.Rename(tmpfile.Name(), tmpfile.Name()+suffix)
	if err != nil {
		return nil, err
	}
	newFile, err := os.Open(tmpfile.Name() + suffix)
	if err != nil {
		return nil, err
	}

	return newFile, nil
}
