package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	httpMethodNotAllowed = "Method Not Allowed"
	httpBadRequest       = "Bad Request"
	httpInternalError    = "Internal Error"
)

var (
	supportedLanguages = []string{"java", "cpp", "python", "javascript"}
	languageSuffix     = map[string]string{
		"java":       "java",
		"python":     "py",
		"cpp":        "cpp",
		"javascript": "js",
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

func checkLanguage(code Code) (in bool) {
	in = false
	for _, language := range supportedLanguages {
		if code.Language == language {
			in = true
			break
		}
	}
	return
}

func lintHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %v", r.RemoteAddr)
	defer r.Body.Close()
	if r.Method != "POST" {
		http.Error(w, httpMethodNotAllowed, 405)
		return
	}

	code := Code{
		Language: strings.ToLower(r.PostFormValue("language")),
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
	if linter == nil {
		log.Fatal("No available linter found.")
		http.Error(w, httpInternalError, 500)
		return
	}

	result, err := linter(tmpfile)
	if err != nil {
		log.Fatal(err)
		http.Error(w, httpInternalError, 500)
		return
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", lintHandler)
	log.Println("server is listenning on port 48722")
	err := http.ListenAndServe(":48722", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
