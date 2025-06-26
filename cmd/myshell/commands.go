package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type ParsedCommand struct {
	Cmd       string
	Args      []string
	RedirType RedirectionType
	RedirFile string
	PipedCmd  *ParsedCommand
}

func handleTypeCmd(args []string) {
	if len(args) == 0 {
		fmt.Println("type: missing argument")
		return
	}

	if slices.Contains(builtInCommands, args[0]) {
		fmt.Printf("%s is a shell builtin\n", args[0])
	} else if path, ok := isOnPath(args[0]); ok {
		fullPath := path + "/" + args[0]
		fmt.Printf("%s is %s\n", args[0], fullPath)
	} else {
		fmt.Printf("%s not found\n", args[0])
	}
}

func handlePwdCmd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working dir : ", err)
		return
	}

	fmt.Println(dir)
}

func handleCdCmd(args []string) {
	path := args[0]

	if path == "~" {
		path = os.Getenv("HOME")
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", path)
		return
	}
}

func displayCmd(cmd ParsedCommand, input *strings.Builder) {
	// clear the input and terminal
	fmt.Print("\r\033[K")
	input.Reset()

	fmt.Printf("$ %s", cmd.Cmd)
	input.WriteString(cmd.Cmd)

	for _, arg := range cmd.Args {
		fmt.Printf(" %s", arg)
		input.WriteString(" " + arg)
	}
}

func handlePipeCmd(cmd *ParsedCommand) {
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		return
	}

	// Execute left CMD
	if slices.Contains(builtInCommands, cmd.Cmd) {
		// For builtin commands, redirect stdout
		originalStdout := os.Stdout
		os.Stdout = writer

		executeBuiltinCommand(cmd)

		os.Stdout = originalStdout
		writer.Close()
	} else {
		// For external commands, use exec with pipe
		leftCmd := exec.Command(cmd.Cmd, cmd.Args...)
		leftCmd.Stdout = writer
		leftCmd.Stderr = os.Stderr
		leftCmd.Stdin = os.Stdin

		go func() {
			defer writer.Close()
			leftCmd.Run()
		}()
	}

	// Execute right CMD with input from left side
	originalStdin := os.Stdin
	os.Stdin = reader

	handleCommand(cmd.PipedCmd) // Recursive for multiple pipes

	os.Stdin = originalStdin
	reader.Close()
}

func executeBuiltinCommand(cmd *ParsedCommand) {
	switch cmd.Cmd {
    case "echo":
        if len(cmd.Args) > 0 {
            fmt.Println(strings.Join(cmd.Args, " "))
        }
    case "type":
        handleTypeCmd(cmd.Args)
    case "pwd":
        handlePwdCmd()
    case "cd":
        handleCdCmd(cmd.Args)
    case "history":
        displayCmdHistory(cmd.Args)
    }
}
