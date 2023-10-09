package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogMessage struct {
	Title   string
	Body    string
	Objects []interface{}
}

func printLog(message LogMessage, typeLog string) {
	f, err := os.OpenFile(fmt.Sprint(time.Now().Format("20060102"), ".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln("Error on close file LogMessage")
		}
	}(f)

	logger := log.New(f, fmt.Sprint(typeLog, "\t"), log.Ldate|log.Ltime)
	logger.Println(message)
}

func (l LogMessage) Error() {
	printLog(l, "ERROR")
}

func (l LogMessage) Info() {
	printLog(l, "INFO")
}
