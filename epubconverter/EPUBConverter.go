package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"gotongwen"
	"log"
	"macgyver"
	"os"
)

type stringHandler func(input string) (output string)

func printLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	convertEPUB("test.epub", "result.epub")
}

func convertEPUB(source string, target string) {

	//Open zip file
	r, err := zip.OpenReader(source)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	//Create result zip file
	targetFile, err := os.Create(target)
	if err != nil {
		log.Fatal(err)
	}
	zipWriter := zip.NewWriter(targetFile)
	defer zipWriter.Close()

	//travel all file in zip
	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		zipFileReader := bufio.NewReader(rc)
		fmt.Println(f.Name)
		if f.Name[len(f.Name)-5:] == "xhtml" {

			//Create a string list from file in zip.
			lines, _ := macgyver.ReadLines_FromReader(zipFileReader)

			//Convert the content by tongweng table
			lines = handleLines(lines, gotongwen.Convert)

			//Write to new zip file
			writeLines_To_ZipFile(lines, zipWriter, f.Name)
		} else {
			writeBytes_To_ZipFile(zipFileReader, zipWriter, f.Name)
		}

		rc.Close()
	}
}

func writeLines_To_ZipFile(lines []string, writer *zip.Writer, fileName string) {
	f, err := writer.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range lines {
		_, err = f.Write([]byte(line + "\n"))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func writeBytes_To_ZipFile(reader *bufio.Reader, writer *zip.Writer, fileName string) {
	var (
		part   []byte
		length int
	)

	part = make([]byte, 1024)
	f, err := writer.Create(fileName)

	for {
		if length, err = reader.Read(part); err != nil {
			break
		}
		fmt.Println(length)
		if length != 0 {
			f.Write(part[:length])
		} else {
			break
		}
	}

	return
}

func handleLines(lines []string, handler stringHandler) (results []string) {
	for _, line := range lines {
		results = append(results, handler(line))
	}

	return
}
