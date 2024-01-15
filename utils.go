package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func logToFile(texts ...any) {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)
	for _, t := range texts {
		log.Print(t)
	}
}

func waitRandomSeconds() {
	// Generate a random number between 1 and 10
	seconds := rand.Intn(10) + 1

	// Wait for the random number of seconds
	time.Sleep(time.Duration(seconds) * time.Second)
}

// convertSizeToHumanRedable convert size in bytes to human readable format
func convertSizeToHumanRedable(size int64) string {
	const (
		KB int64 = 1024
		MB int64 = KB * 1024
		GB int64 = MB * 1024
	)

	var sizeReadable string

	if size >= GB {
		sizeReadable = fmt.Sprintf("%.2fGB", float64(size)/float64(GB))
	} else if size > 1000000 {
		sizeReadable = fmt.Sprintf("%.2fMB", float64(size)/float64(MB))
	} else if size > 1000 {
		sizeReadable = fmt.Sprintf("%.2fKB", float64(size)/float64(KB))
	} else {
		sizeReadable = fmt.Sprintf("%dBytes", size)
	}
	return sizeReadable

}
