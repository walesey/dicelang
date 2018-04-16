package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go/build"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

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

	codeCache := new(sync.Map)
	http.HandleFunc("/code", func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Error: ", err)))
			return
		}

		var codes []string
		if err = json.Unmarshal(body, &codes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprint("Error: ", err)))
			return
		}

		results := make([]interface{}, len(codes))
		for i, code := range codes {
			if cachedResult, ok := codeCache.Load(code); ok {
				results[i] = cachedResult
				continue
			}

			log.Println("Parsing: ", code)
			buf := bytes.NewReader([]byte(code))
			var output interface{}
			if output, err = parser.NewParser(buf).Execute(); err != nil {
				log.Println("Error Parsing: ", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprint("Parse Error: ", err)))
				return
			} else {
				results[i] = output
			}
			codeCache.Store(code, output)
		}

		if output, err := json.Marshal(results); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint("Json Marshal Error: ", err)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(output)
		}
	})

	// periodically flush the code cache
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			codeCache.Range(func(key, value interface{}) bool {
				codeCache.Delete(key)
				return true
			})
		}
	}()

	log.Fatal(http.ListenAndServe(":4000", nil))
}
