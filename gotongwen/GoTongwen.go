package gotongwen

import (
	"bytes"
	"strings"
)

func Convert(str string) (result string) {
	//convert word
	str = convertWord(str)
	//convert phrase
	result = convertPhrase(str)
	return
}

func convertWord(str string) (result string) {
	//convert word
	strArray := strings.Split(str, "")
	for i := 0; i < len(strArray); i++ {
		value, ok := s2tTable[strArray[i]]
		if ok {
			strArray[i] = value
		}
	}
	result = strings.Join(strArray, "")
	return
}

func convertPhrase(str string) (result string) {
	//convert phrase
	buffer := bytes.NewBufferString("")
	strLen := len(str)

	//travel all word and find match phrase
	for index := 0; index < strLen; {

		//Avoid out of ranage
		subStrLength := minLength(maxLength, strLen-index)

		//If there is match phrase , replace by it
		if matchPhrase, isMatch := findMatchPhrase(str[index:index+subStrLength], subStrLength); isMatch {
			buffer.WriteString(matchPhrase)
			index += len(matchPhrase)
		} else {
			buffer.WriteString(str[index : index+1])
			index++
		}
	}

	result = buffer.String()
	return
}

func findMatchPhrase(subStr string, maxPhraseLength int) (matchPhrase string, isMatch bool) {
	//find the match phrase and replace string by it
	for phraseIndex := maxPhraseLength; phraseIndex > 1; phraseIndex-- {

		value, ok := s2tTable[subStr[:phraseIndex]]

		if ok {
			matchPhrase = value
			isMatch = true
			return
		}
	}

	isMatch = false
	return
}

func minLength(len1 int, len2 int) (result int) {
	result = len1

	if len1 > len2 {
		result = len2
	}
	return
}
