package knife

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func WriteLines(path string, lines []string) {
	var file *os.File
	var err error

	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(item + "\n")
		if err != nil {
			log.Fatal(err)
			break
		}
	}

	return
}

func ReadLines(path string) (lines []string, err error) {
	var (
		file *os.File
	)

	//Open a file
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	lines, err = ReadLines_FromFile(file)

	return
}

func ReadLines_FromFile(file *os.File) (lines []string, err error) {
	var (
		part   []byte
		prefix bool
	)
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

func ReadLines_FromReader(reader *bufio.Reader) (lines []string, err error) {
	var (
		part   []byte
		prefix bool
	)
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
