package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtInCommands = []string{"exit", "echo", "type", "pwd"}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading input:", err)
			continue
		}

		userInput = strings.ReplaceAll(userInput, "\n", "")
		cmd, args := parseInput(userInput)

		handleCommand(cmd, args)
	}
}

func handleCommand(cmd string, args []string) {
	switch cmd {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "type":
		if len(args) == 0 {
			fmt.Println("type: missing argument")
			return
		}

		if listContains(builtInCommands, args[0]) {
			fmt.Printf("%s is a shell builtin\n", args[0])
		} else if path, ok := isOnPath(args[0]); ok {
			fullPath := path + "/" + args[0]
			fmt.Printf("%s is %s\n", args[0], fullPath)
		} else {
			fmt.Printf("%s not found\n", args[0])
		}
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting working dir : ", err)
			return
		}

		fmt.Println(dir)
	case "cd":
		path := args[0]

		if path == "~" {
			path = os.Getenv("HOME")
		}

		err := os.Chdir(path)
		if err != nil {
			fmt.Printf("cd: %s: No such file or directory\n", path)
			return
		}
	default:
		runCommand(cmd, args)
	}
}

func parseInput(input string) (string, []string) {
	fields := strings.Fields(input)
	cmd := fields[0]

	if len(fields) > 1 {
		return cmd, fields[1:]
	} else {
		return cmd, nil
	}
}

func listContains(list []string, elem string) bool {
	for _, e := range list {
		if e == elem {
			return true
		}
	}

	return false
}

func isOnPath(cmd string) (string, bool) {
	osPath := os.Getenv("PATH")
	paths := strings.Split(osPath, ":")

	for _, path := range paths {
		if _, err := os.Stat(path + "/" + cmd); err == nil {
			return path, true
		}
	}

	return "", false
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
