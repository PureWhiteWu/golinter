package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type lintFunc func(*os.File) (Result, error)

func lintJava(tmpfile *os.File) (Result, error) {
	cmd := exec.Command("java", "-jar", "linters/java/checkstyle-7.7-all.jar", "-c", "linters/java/sun_checks.xml", tmpfile.Name())
	b, _ := cmd.CombinedOutput()
	re, _ := regexp.Compile("(\\[ERROR].*?)\n")
	tmp := re.FindAllString(string(b), -1)
	var s []string
	// replace unnecessary information with ""
	for _, str := range tmp {
		re, _ = regexp.Compile("\n")
		str = re.ReplaceAllString(str, "")
		re, _ = regexp.Compile("/.*?java:")
		str = re.ReplaceAllString(str, "")
		s = append(s, str)
	}

	result := Result{
		ErrorNum: len(s),
		Errors:   s,
	}
	return result, nil
}

// dispatch linters
func dispatch(code Code) lintFunc {
	switch code.Language {
	case "java":
		fallthrough
	case "Java":
		return lintJava
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
	suffix := "." + languageSuffix[strings.ToLower(code.Language)]
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
