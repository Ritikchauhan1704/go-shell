package main

import (
	"bufio"
	"fmt"
	"os"
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

		argv := shell.SplitCmd(command)
		cmd := argv[0]

		writeOutput := false
		writeError := false
		mode := 'a'
		var file *os.File
		origStdout := os.Stdout
		origStderr := os.Stderr
		for ind, arg := range argv {
			if ind+1 < len(argv) && (arg == ">" || arg == "1>" || arg == ">>" || arg == "1>>") {
				if arg == ">" || arg == "1>" {
					mode = 'w'
				}
				writeOutput = true
				file, err = shell.CreateFile(argv[len(argv)-1], mode)
				argv = argv[:ind]

			} else if ind+1 < len(argv) && (arg == "2>" || arg == "2>>") {
				if arg == "2>" {
					mode = 'w'
				}
				writeError = true
				file, err = shell.CreateFile(argv[len(argv)-1], mode)
				argv = argv[:ind]

			}
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if writeOutput {
			os.Stdout = file
			defer file.Close()

		}
		if writeError {
			os.Stderr = file
			defer file.Close()

		}

		switch cmd {
		case "type":
			shell.TypFun(argv)
		case "echo":
			shell.EchoCmd(argv)
		case "exit":
			shell.ExitCmd(argv)
		case "pwd":
			shell.Pwd(argv[len(argv)-1])
		case "cd":
			shell.Cd(argv)
		default:
			shell.ExtProg(argv)

		}

		os.Stdout = origStdout
		os.Stderr = origStderr
	}

}
