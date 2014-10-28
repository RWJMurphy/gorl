package main

import (
	"github.com/RWJMurphy/gorl/lib"
)

func main() {
	g, err := gorl.NewGame()
	if err != nil {
		panic(err)
	}
	defer g.Close()
	g.Run()
}
