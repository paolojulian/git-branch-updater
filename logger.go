package main

import (
	"fmt"
	"log"
	"strings"
)

type Logger interface {
	Header(number int, title string)
	Description(description string)
	Error(err error, description ...string)
}

type LoggerImpl struct {
}

func NewLogger() *LoggerImpl {
	return &LoggerImpl{}
}

func (l *LoggerImpl) Header(number int, title string) {
	fmt.Printf("\n-- %d. %s", number, title)
}

func (l *LoggerImpl) Description(description string) {
	fmt.Printf("\n---- %s", description)
}

func (l *LoggerImpl) Error(err error, description ...string) {
	fmt.Println("\n=======================================")
	log.Fatalln(`
ERROR: ` + err.Error() + `
Description: ` + strings.Join(description, " ") + `
=======================================
	`)
}
