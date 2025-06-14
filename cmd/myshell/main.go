package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtInCommands = []string{"exit", "echo", "type", "pwd", "cd"}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("error reading input:", err)
			os.Exit(1)
		}

		cmd, args := parseInput(userInput[:len(userInput)-1])

		handleCommand(cmd, args)
	}
}

func handleCommand(cmd string, args []string) {
	switch cmd {
	case "exit":
		os.Exit(0)
	case "echo":
		if len(args) > 0 {
			fmt.Println(strings.Join(args, " "))
		}
	case "type":
		handleTypeCmd(args)
	case "pwd":
		handlePwdCmd()
	case "cd":
		handleCdCmd(args)
	default:
		runCommand(cmd, args)
	}
}

func runCommand(cmd string, args []string) {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		fmt.Printf("%s: command not found\n", cmd)
	}
}
