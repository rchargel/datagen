package main

import (
	"math/rand"
	"strconv"
	"strings"
)

func (s staticFile) generateData(channel chan string) {
	l := len(s.possibleValues)
	num := int(s.numberOfEntities)
	for entityID := 1; entityID <= num; entityID++ {
		channel <- strings.Join(s.randValue(entityID, l), ",")
	}
	close(channel)
}

func (s staticFile) randValue(entityID int, l int) []string {
	return []string{strconv.Itoa(entityID), s.possibleValues[rand.Intn(l)]}
}

func (s staticFile) isStatic() bool {
	return true
}

func (s staticFile) getFilePath() string {
	return s.filePath
}
