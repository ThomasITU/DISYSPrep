package helpermethod

import (
	"log"
	"os"
)

// helper method to help find error locations
func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("happened inside method: %s err: %v", msg, err)
	}
}

// log message in file
func Logger(message string, logFileName string) {
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(message)
}
