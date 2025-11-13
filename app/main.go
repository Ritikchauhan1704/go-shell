package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
"github.com/ShwetaRoy17/go-shell/app/shell" 
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout




func main() {
	for true {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		argv := strings.Fields(command)
		cmd := argv[0]

		switch cmd {
		case "type":
			shell.TypFun(argv)
		case "echo":
			shell.EchoCmd(argv)
		case "exit":
			shell.ExitCmd(argv)
		case "pwd":
			shell.Pwd()
		case "cd":
			shell.Cd(argv)
		default:
			shell.ExtProg(argv)


		}

	}

}


