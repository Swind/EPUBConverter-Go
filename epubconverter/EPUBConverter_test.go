package main

import (
	"runtime"
	"testing"
)

func BenchmarkConvertCore4(b *testing.B) {
	executeConvert("./testdata", 4)
}

func executeConvert(sourcePath string, coreSize int) {

	runtime.GOMAXPROCS(coreSize)

	request_queue = make(chan int, coreSize)

	travelAllFile(sourcePath)
}
