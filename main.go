package main

import (
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		if strings.HasPrefix(err.Error(), "exit status") {
			os.Exit(0)
		}
		panic(err)
	}
}

func main() {
	switch os.Args[1] {
	case "run":
		run(os.Args[2:]...)
	default:
		panic("no such command")
	}
}

func run(args ...string) {

}

