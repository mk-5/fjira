package app

import (
	"log"
	"github.com/atotto/clipboard"
)

func CopyIssue(issue string) {
	err := clipboard.WriteAll(issue)
	if err != nil {
		log.Fatal(err)
	}
}
