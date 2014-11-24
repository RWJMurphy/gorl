package gorl

import (
	"log"
	"math/rand"
	"os"
	"time"
)

const logFilePath = "gorl.log"

type GorlCLI interface {
	Run()
	Close()
}

type gorlCLI struct {
	logFile *os.File
	log     *log.Logger
	game    *Game
}

func NewCLI(args []string) GorlCLI {
	cli := gorlCLI{}
	logFile, err := os.OpenFile(
		logFilePath,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0666,
	)
	if err != nil {
		panic(err)
	}
	cli.logFile = logFile
	cli.log = log.New(logFile, "gorl: ", log.Ldate|log.Ltime|log.Lshortfile)
	cli.log.Println("Starting gorl")
	seed := time.Now().UnixNano()
	dice := rand.New(rand.NewSource(seed))
	game, err := NewGame(cli.log, dice)
	if err != nil {
		cli.log.Panic(err)
	}
	cli.game = game
	return &cli
}

func (cli *gorlCLI) Run() {
	cli.game.Run()
}

func (cli *gorlCLI) Close() {
	cli.logFile.Sync()
	cli.logFile.Close()
	cli.game.Close()
}
