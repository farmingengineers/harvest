package main

import (
	"io"
	"log"
	"os"
)

var l = func() *log.Logger {
	w := io.Discard
	if os.Getenv("DEBUG") != "" {
		f, err := os.Create("debug.log")
		if err != nil {
			panic(err)
		}
		w = f
	}
	l := log.New(w, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	l.Printf("started")
	return l
}()
