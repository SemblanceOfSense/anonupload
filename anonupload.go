package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileList struct {
	Files []string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			body, err := os.ReadFile("/var/www/anonupload/index.html")
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

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			content, err := jpeg.Decode(r.Body)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rand := rand.Int()
			f, err := os.Create("/var/www/anonupload/files/" + filepath.Base(strconv.Itoa(rand) + ".jpeg"))
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			err = jpeg.Encode(f, content, &jpeg.Options{
				Quality: 100,
			})
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
				if (strings.Contains(v.Name(), ".jpeg")) {
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
