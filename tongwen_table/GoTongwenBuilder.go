package tongwen_table

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func CreateTongwenTable() {

	word_s2t, _ := readLines("word_s2t.txt")

	phrase_s2t, _ := readLines("phrase_s2t.txt")

	mergeList := make([]string, len(word_s2t)+len(phrase_s2t))
	copy(mergeList, word_s2t)
	copy(mergeList[len(word_s2t):], phrase_s2t)

	resultList := createTongWenMapString(mergeList)

	writeLines("GoTongwenTable.go", resultList)

	return
}

func createTongWenMapString(lines []string) []string {
	var resultLines []string

	resultLines = append(resultLines, "package gotongwen")
	resultLines = append(resultLines, "var s2tTable = map[string] string {")

	maxLength := 0

	//Remove duplicate key
	set := make(map[string]string)
	for _, line := range lines {
		arr := strings.Split(line, ",")
		set[arr[0]] = arr[1]
	}

	for key, value := range set {
		if len(key) > maxLength {
			maxLength = len(key)
		}

		resultLines = append(resultLines, fmt.Sprintf("\"%s\": \"%s\",", key, value))
	}

	resultLines = append(resultLines, "}")

	resultLines = append(resultLines, fmt.Sprintf("var maxLength =%d", maxLength))

	return resultLines
}

func writeLines(path string, lines []string) {
	var file *os.File
	var err error

	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(item + "\n")
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	return
}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)

	//Open a file
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))

	//Read line from file and append the line to the List
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}

		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}
