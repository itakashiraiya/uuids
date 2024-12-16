package main

import (
	"fmt"
	"os"
	"time"
	// "github.com/itakashiraiya/uuids/internals/uuids"
	"golang.org/x/term"
)

func main() {
	// uuids.Test()
	a := []string{
		"test1.1",
		"test1.2",
		"test1.3",
	}

	render(a)

	time.Sleep(time.Second * 1)

	a = []string{
		"test2.1",
		"test2.1",
		"test2.1",
	}

	render(a)
	for range a {
		fmt.Println()
	}
}

func loop() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")
	defer fmt.Print("\r\033[K")
}

func internal() {

}

func render(outputs []string) {
	for _, v := range outputs {
		fmt.Printf("\r\033[K %s\n", v)
	}
	fmt.Printf("\033[%dA", len(outputs))
}
