package main

import (
	"log"
	"os"
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
