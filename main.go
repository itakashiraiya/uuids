package main

import (
	"errors"
	"fmt"
	"math/big"

	// "math/rand"
	"os"

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

type inputCmd string

func loop() {
	// rand := rand.New(rand.NewSource(4))
	pos := big.NewInt(0)
	usrString := ""
	internal(usrString, pos)
	for {
		// pos := new(big.Int).Rand(rand, new(big.Int).Add(uuids.MaxEntropyNum(), big.NewInt(1)))
		input, err := input()
		if err != nil {
			fmt.Println(err.Error())
		}
		if val, ok := input.(inputCmd); ok {
			switch val {
			case "A":
				if pos.Cmp(big.NewInt(0)) > 0 {
					pos.Add(pos, big.NewInt(-1))
				}
			case "B":
				if pos.Cmp(uuids.MaxEntropyNum()) < 0 {
					pos.Add(pos, big.NewInt(1))
				}
			case "exit":
				return
			}
		} else if val, ok := input.(byte); ok {
			usrString += string(val)
		}
		internal(usrString, pos)
	}
}

func input() (interface{}, error) {
	var input [3]byte
	n, err := os.Stdin.Read(input[:])
	if err != nil {
		panic(err)
	}
	if n == 1 {
		return input[0], nil
	}

	if n == 2 && input[0] == 27 && input[1] == 27 {
		return inputCmd("exit"), nil
	}

	if n == 3 && input[0] == 27 {
		if input[1] == 91 {
			switch input[2] {
			case 'A', 'B':
				return inputCmd(input[2]), nil
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("unknown input (n = %d, bytes = %d\r", n, input))
}

func inputOld() (escSeq bool, input [3]byte, err error) {
	n, err := os.Stdin.Read(input[:])

	if n == 3 && input[0] == 27 && input[1] == 91 {
		return true, input, nil
	}
	return false, input, nil
}

func internal(usrString string, pos *big.Int) {
	display := []string{usrString}
	uuids := uuids.GetUuids(pos, display_height)
	display = append(display, uuids...)
	render(display)
}

func render(outputs []string) {
	for _, v := range outputs {
		fmt.Printf("\r\033[K %s\n", v)
	}
	fmt.Printf("\033[%dA", len(outputs))
}
