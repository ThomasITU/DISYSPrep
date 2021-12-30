package HelperMethod

import (
	"log"
	"os"
)

const (
	FRONT_END_ADDRESS = "localhost:5000"
	MAX_REPLICAS = 5
)

type Value struct {
	Value  int64
	UserId int64
}

// helper method to help find error locations
func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("happened inside method: %s err: %v", msg, err)
	}
}

// log message in file
func Logger(message string, logFileName string) {
	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(message)
}
