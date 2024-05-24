package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		cmd = strings.ReplaceAll(cmd, "\n", "")
		handleCommand(cmd)
	}
}

func handleCommand(cmd string) {
	fmt.Printf("%s: command not found\n", cmd)
}
