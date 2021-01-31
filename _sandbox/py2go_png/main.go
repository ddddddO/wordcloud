package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(buf)

	io.Copy(os.Stdout, r)
}
