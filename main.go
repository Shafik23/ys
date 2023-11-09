package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/shafik23/ys/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("You are Wise %s ... \n", user.Username)

	repl.Start(os.Stdin, os.Stdout)
}
