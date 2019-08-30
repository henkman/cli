package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/inancgumus/screen"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("obs delay command")
		return
	}
	delay, err := time.ParseDuration(os.Args[1])
	if err != nil {
		panic(err)
	}

	screen.Clear()
	for {
		screen.MoveTopLeft()
		cmd := exec.Command(os.Args[2], os.Args[3:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
		time.Sleep(delay)
	}
}
