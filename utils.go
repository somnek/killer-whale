package main

import (
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
