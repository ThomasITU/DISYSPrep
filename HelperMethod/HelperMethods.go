package helperMethod

import "log"

// helper method to help find error locations
func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("happened inside method: %s err: %v", msg, err)
	}
}
