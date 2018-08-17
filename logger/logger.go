package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	enabledFileLog = true
	logFilename    = "logg.txt"
	mu             sync.Mutex
	filter         string
)

func Log(text string) {
	if filter == "" || (filter != "" && strings.Contains(text, filter)) {
		log.Print(text)
		if !enabledFileLog {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if err := writeToFile(logFilename, text); err != nil {
			log.Println(err)
		}
	}
}
func Logln(v ...interface{}) {
	Log(fmt.Sprintln(v...))
}

func Logf(format string, v ...interface{}) {
	Log(fmt.Sprintf(format, v...))
}

func SetFilename(filename string) {
	logFilename = filename
}

func Disable() {
	enabledFileLog = false
}

func Enable() {
	enabledFileLog = true
}

func Filter(f string) {
	filter = f
}
func writeToFile(fileName, text string) error {
	file, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	text = fmt.Sprintf("%s %s", time.Now().String(), text)
	if _, err = file.WriteString(text); err != nil {
		return err
	}
	return nil
}
