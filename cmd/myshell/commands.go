package main

import (
	"fmt"
	"os"
	"slices"
)

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
