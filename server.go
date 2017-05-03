package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const (
	httpMethodNotAllowed = "Method Not Allowed"
	httpBadRequest       = "Bad Request"
	httpInternalError    = "Internal Error"
)

var (
	supportedLanguages = []string{"java", "Java"}
	languageSuffix     = map[string]string{
		"java":   "java",
		"python": "py",
		"cpp":    "cpp",
	}
)

type Code struct {
	Language string `json:"language"`
	Source   string `json:"source"`
}

type Result struct {
	ErrorNum int      `json:"error_num"`
	Errors   []string `json:"errors"`
}

func checkLanguage(code Code) bool {
	in := false
	for _, language := range supportedLanguages {
		if code.Language == language {
			in = true
			break
		}
	}
	if in {
		return true
	}
	return false
}

func lintHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "POST" {
		http.Error(w, httpMethodNotAllowed, 405)
		return
	}

	code := Code{
		Language: r.PostFormValue("language"),
		Source:   r.PostFormValue("source"),
	}
	if !checkLanguage(code) {
		http.Error(w, httpBadRequest, 400)
		return
	}

	tmpfile, err := createTempFile(code)
	if err != nil {
		log.Fatal(err)
		http.Error(w, httpInternalError, 500)
		return
	}

	defer func() {
		err := os.Remove(tmpfile.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer tmpfile.Close()
	linter := dispatch(code)
	result, err := linter(tmpfile)
	if err != nil {
		log.Fatal(err)
		http.Error(w, httpInternalError, 500)
		return
	}
	if err = json.NewEncoder(w).Encode(result); err {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", lintHandler)
	err := http.ListenAndServe(":48722", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
