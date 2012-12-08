package gotongwen

import (
	"testing"
)

func Test_Convert_String(t *testing.T) {
	simple := "HTTP 2.0 发布了首个草案，该草案直接复制于"
	expectResult := "HTTP 2.0 發佈了首個草案，該草案直接複製於"

	result := Convert(simple)

	if expectResult != result {
		t.Error("Before=", simple)
		t.Error("After =", result)
		t.Error("Expect=", expectResult)
	} else {
		t.Log("Result =", result)
	}
}
