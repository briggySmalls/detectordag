package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	suffixOfInterest = ".html"
	outputName       = "html_data.go"
	packageName      = "email"
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
	out.Write([]byte(fmt.Sprintf("package %s \n\nconst (\n", packageName)))
	// Scan the directory
	fs, _ := ioutil.ReadDir(dir)
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), suffixOfInterest) {
			// File is an html file
			log.Print(f.Name())
			var err error
			// Add a variable with the same name as the file
			_, err = out.Write([]byte(
				fmt.Sprintf(
					"%sHtmlTemplateSource = `",
					strings.TrimSuffix(f.Name(), suffixOfInterest),
				),
			))
			// Read the file, and copy it into our output
			f, err := os.Open(filepath.Join(dir, f.Name()))
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(out, f)
			if err != nil {
				log.Fatal(err)
			}
			f.Close()
			// Terminate the line
			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}
