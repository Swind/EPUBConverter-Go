package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"gotongwen"
	"knife"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const MAX_BUFFLENGTH int = 1024 * 16
const sourcePath string = "./source/"
const targetPath string = "./target/"

type stringHandler func(input string) (output string)

var request_queue chan int
var count int = 0

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	request_queue = make(chan int, runtime.NumCPU())

	fmt.Printf("CPU is %d core\n", runtime.NumCPU())
	err := travelAllFile(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

}

func travelAllFile(path string) (err error) {

	err = filepath.Walk(path, folderVisit)
	if err != nil {
		log.Fatal(err)
	}
	return

}

func folderVisit(path string, f os.FileInfo, err error) error {

	if strings.HasSuffix(f.Name(), ".epub") {
		fmt.Printf("Convert: %s\n", gotongwen.Convert(f.Name()))

		//Put a request to queue, the max size of queue is CPU number.
		count++
		request_queue <- count
		go convertEPUB(path, targetPath+gotongwen.Convert(f.Name()))
	}

	return nil
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

	//Travel all file in zip
	for _, f := range r.File {

		//Open file in the zip file and create a buffer reader
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		zipFileReader := bufio.NewReader(rc)

		//We only convert xhtml file , other files will be wrote directly.
		if endWith(f.Name, "xhtml") || endWith(f.Name, "html") || endWith(f.Name, "htm") {

			//Create a string list from file in zip.
			lines, _ := knife.ReadLines_FromReader(zipFileReader)

			//Convert the content by tongweng table
			lines = handleLines(lines, gotongwen.Convert)

			//Write to new zip file
			writeLines_To_ZipFile(lines, zipWriter, f.Name)

		} else {

			writeBytes_To_ZipFile(zipFileReader, zipWriter, f.Name)

		}

		rc.Close()
	}

	<-request_queue
}

func endWith(source string, suffix string) bool {
	return strings.HasSuffix(source, suffix)
}

func writeLines_To_ZipFile(lines []string, writer *zip.Writer, fileName string) {

	//Create a file in the zip file
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

	var part = make([]byte, 1024)
	var length int = 0

	f, err := writer.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	for {
		if length, err = reader.Read(part); err != nil {
			break
		}

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
