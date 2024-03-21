package utils

import (
	"log"
	"strings"
)

// LogErrorAndContinue is a helper to check if an error happened, log it, and continue
func LogErrorAndContinue(err error, context ...string) {
	combinedContext := strings.Join(context, " ")
	if err != nil {
		log.Printf("ERROR: %s %s\n", combinedContext, err.Error())
	}
}
