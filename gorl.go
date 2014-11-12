package main

import (
	"log"
	"os"
	// "runtime"

	"github.com/RWJMurphy/gorl/lib"
)

const logFilePath = "gorl.log"

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	logFile, err := os.OpenFile(
		logFilePath,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0666,
	)
	if err != nil {
		panic(err)
	}
	defer logFile.Sync()
	defer logFile.Close()
	log := log.New(logFile, "gorl: ", log.Ldate|log.Ltime|log.Lshortfile)
	log.Println("Starting gorl")
	g, err := gorl.NewGame(log)
	if err != nil {
		log.Panic(err)
	}
	defer g.Close()
	g.Run()
}
