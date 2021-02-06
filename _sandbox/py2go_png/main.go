package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var pngByte []byte

func main() {
	args := os.Args
	if len(args) == 1 {
		io.Copy(os.Stdout, os.Stdin)
		return
	}

	if args[1] == "serve" {
		log.Println("start")

		var err error
		pngByte, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		http.HandleFunc("/", pngHandler)

		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}
}

func pngHandler(w http.ResponseWriter, r *http.Request) {
	reader := bytes.NewReader(pngByte)
	img, _, err := image.Decode(reader)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "image/png")
	err = png.Encode(w, img)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
