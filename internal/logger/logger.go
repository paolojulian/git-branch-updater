package logger

import (
	"fmt"
	"log"
	"strings"
)

var CURRENT_HEADER int = 1

type Logger interface {
	Header(title string)
	Description(description string)
	Error(err error, description ...string)
}

type LoggerImpl struct {
}

func NewLogger() *LoggerImpl {
	return &LoggerImpl{}
}

func (l *LoggerImpl) Header(title string) {
	fmt.Printf("\n-- %d. %s", CURRENT_HEADER, title)
	CURRENT_HEADER++
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
