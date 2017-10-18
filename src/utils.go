package main

import (
	"os"
	"strings"
)

// Turn a list of errors into a list of strings that can easily be printed
func GetErrorMessages(errs []error) []string {
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = err.Error()
	}
	return messages
}

// Print a message out to standard error and exit the program
func Die(message string) {
	os.Stderr.WriteString(strings.Join([]string{message, "\n"}, ""))
	os.Exit(2)
}
