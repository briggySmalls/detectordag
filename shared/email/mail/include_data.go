package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	suffixOfInterest = ".html"
	outputName       = "html_data.go"
)

// Reads all .txt files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	// Grab the directory to work on
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		panic("Need input directory")
	}
	dir := argsWithoutProg[0]
	// Create the file we'll dump out results into
	out, _ := os.Create(outputName)
	defer out.Close()
	out.Write([]byte("package email \n\nconst (\n"))
	// Scan the directory
	fs, _ := ioutil.ReadDir(dir)
	for _, f := range fs {
		log.Print(f.Name())
		if strings.HasSuffix(f.Name(), suffixOfInterest) {
			// File is an html file
			// Add a variable with the same name as the file
			out.Write([]byte(strings.TrimSuffix(f.Name(), suffixOfInterest) + " = `"))
			// Read the file, and copy it into our output
			f, _ := os.Open(f.Name())
			io.Copy(out, f)
			f.Close()
			// Terminate the line
			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}
