package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go/build"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/walesey/dicelang/parser"
)

func init() {
	// set working dir to access assets
	p, _ := build.Import("github.com/walesey/dicelang", "", build.FindOnly)
	os.Chdir(p.Dir)
}

func main() {
	fs := http.FileServer(http.Dir("./website/app"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		codeB64 := req.URL.Query().Get("code")
		code, err := base64.URLEncoding.DecodeString(codeB64)
		if len(code) == 0 || err != nil {
			code = []byte("hist 2d6.add")
		}

		t, err := template.ParseFiles("./website/app/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error Parsing Template: ", err)))
		}
		if err := t.Execute(w, map[string]string{"code": string(code)}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error Executing Template: ", err)))
		}
	})

	http.HandleFunc("/code", func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error: ", err)))
		}

		log.Println("Parsing: ", string(body))
		buf := bytes.NewReader(body)
		if output, err := parser.NewParser(buf).Execute(); err != nil {
			log.Println("Error Parsing: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Parse Error: ", err)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(output))
		}
	})
	log.Fatal(http.ListenAndServe(":4000", nil))
}
