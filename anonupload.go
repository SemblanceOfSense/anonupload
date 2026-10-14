package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

type FileList struct {
	Files []string
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("/var/www/anonupload/static")))

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			content, err := webp.Decode(r.Body, &decoder.Options{})
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rand := rand.Int()
			f, err := os.Create("/var/www/anonupload/files/" + filepath.Base(strconv.Itoa(rand) + ".webp"))
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			err = webp.Encode(f, content, &encoder.Options{})
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/listfiles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			files, err := os.ReadDir("/var/www/anonupload/files")
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var filelist FileList
			for _, v := range files {
				if (strings.Contains(v.Name(), ".webp")) {
					filelist.Files = append(filelist.Files, v.Name())
				}
			}

			body, err := json.Marshal(filelist)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			status, err := w.Write(body)
			if err != nil {
				fmt.Println(err)
			}
			if status != 200 {
				fmt.Println(http.StatusText(status))
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/getfile", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			queries := r.URL.Query()
			filename := ""
			if queries["filename"] != nil {
				filename = "/var/www/anonupload/files/" + filepath.Base(queries["filename"][0])
			} else {
				fmt.Println("No filename provided")
				http.Error(w, "No filename provided", http.StatusNotFound)
				return
			}
			body, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			status, err := w.Write(body)
			if err != nil {
				fmt.Println(err)
			}
			if status != 200 {
				fmt.Println(http.StatusText(status))
			}
		}
	})

	http.ListenAndServe("localhost:8080", nil)
}
