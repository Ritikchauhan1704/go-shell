package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"os/exec"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

var builtIns = []string{"type", "echo", "exit"}



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
			typFun(argv)
		case "echo":
			EchoCmd(argv)
		case "exit":
			ExitCmd(argv)
		default:
			extProg(argv)


		}

	}

}

func extProg(argv []string) {
	path, exists := FindInPath(argv[0])
	if exists  {
		cmd := exec.Command(path, argv[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error executing command:", err)
		}
	
	} else {
		fmt.Printf("%s: command not found\n", argv[0])
	}

}

func typFun(argv []string) {

	if len(argv) == 1 {
		return
	}

	val := argv[1]

	if slices.Contains(builtIns, val) {
		fmt.Printf("%s is a shell builtin\n",val)
		return
	}
	if file, exists := FindInPath(val); exists == true {
		fmt.Printf("%s is %s\n", val, file)
		return
	}
	fmt.Printf("%s: not found\n", val)

}

func FindInPath(bin string) (string, bool) {
	// this will tell if the command exists in the path or not
    path, err := exec.LookPath(bin)
    if err != nil {
        return "", false
    }
    return path, true
}


func ExitCmd(argv []string) {
	code := 0
	if len(argv) > 1 {
		val, err := strconv.Atoi(argv[1])
		if err != nil {
			code = val
		}
	}
	os.Exit(code)
}

func EchoCmd(argv []string) {
	output := strings.Join(argv[1:], " ")
	fmt.Println(output)
}
