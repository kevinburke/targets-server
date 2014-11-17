package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

// via http://stackoverflow.com/a/22892986/329700
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type response struct {
	Status  int  `json:"status"`
	Success bool `json:"success"`
}

func metricsHandler(w http.ResponseWriter, r *http.Request, directory string) {
	randFile := randSeq(10)
	f, err := os.Create(filepath.Join(directory, randFile))
	if err != nil {
		log.Fatal(err.Error())
	}
	err = os.Chmod(f.Name(), 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "sark/0.1")
	resp := response{
		Status:  200,
		Success: true,
	}
	js, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Write(js)
}

func main() {
	directory := flag.String("directory", "", "Directory to store metrics data in")
	flag.Parse()
	if len(flag.Args()) != 0 || *directory == "" {
		log.Fatal("Please supply a directory")
	}
	err := os.MkdirAll(*directory, 0755)
	if err != nil {
		log.Fatal(err.Error())
	}
	http.HandleFunc("/api/targets/v1/metrics", func(w http.ResponseWriter, r *http.Request) {
		metricsHandler(w, r, *directory)
	})
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
