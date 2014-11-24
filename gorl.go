package main

import (
	"github.com/RWJMurphy/gorl/lib"
	"os"
)

func main() {
	cli := gorl.NewCLI(os.Args[1:])
	defer cli.Close()
	cli.Run()
}
