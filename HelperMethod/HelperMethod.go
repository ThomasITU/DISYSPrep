package helpermethod

import (
	"fmt"
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
	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(message)
}

// for lamport time stamp
func Max(clientTimeStamp int64, serverTimeStamp int64) int64{
	
	if clientTimeStamp <= serverTimeStamp {
		return serverTimeStamp
	}  
	return clientTimeStamp
}
	

// for formating easier to verify
func LoggerWithTimestamp(message string, logFileName string, timestamp int64) {
	msg := fmt.Sprintf("Timestamp(%v): %v", timestamp, message) 
	Logger(msg,logFileName)
}
	

