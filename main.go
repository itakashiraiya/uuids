package main

import (
	"fmt"
	"math/big"
	// "math/rand"
	"os"
	"time"

	"github.com/itakashiraiya/uuids/internals/uuids"
	"golang.org/x/term"
)

var display_height int

func init() {
	display_height = 10
}

func main() {
	start()
}

func start() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")
	defer fmt.Print("\r\033[K")
	defer fmt.Print("\033[999;1H")
	loop()
}

func loop() {
	// rand := rand.New(rand.NewSource(4))
	pos := big.NewInt(0)
	for i := range 1 {
		// pos := new(big.Int).Rand(rand, new(big.Int).Add(uuids.MaxEntropyNum(), big.NewInt(1)))
		internal(new(big.Int).Add(pos, big.NewInt(int64(i))))
		time.Sleep(time.Second * 1)
	}
}

func internal(pos *big.Int) {
	uuids := uuids.GetUuids(pos, display_height)
	render(uuids)
}

func render(outputs []string) {
	for _, v := range outputs {
		fmt.Printf("\r\033[K %s\n", v)
	}
	fmt.Printf("\033[%dA", len(outputs))
}
