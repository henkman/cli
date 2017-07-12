package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: command | pipex program")
		return
	}
	prog, err := exec.LookPath(os.Args[1])
	if err != nil {
		fmt.Println("program not found:", err)
		return
	}
	bin := bufio.NewReader(os.Stdin)
	for {
		line, _ := bin.ReadString('\n')
		if line == "" {
			break
		}
		line = strings.TrimSpace(line)
		args := strings.Split(line, " ")
		cmd := exec.Command(prog, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
