package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var builtIns = map[string]bool{"type": true, "echo": true, "exit": true, "pwd": true}
var extCmd = map[string]bool{"cat": true}

func TypFun(argv []string) {

	if len(argv) == 1 {
		return
	}

	val := argv[1]

	if builtIns[val] {
		fmt.Printf("%s is a shell builtin\n", val)
		return
	}
	if file, exists := FindInPath(val); exists {
		fmt.Printf("%s is %s\n", val, file)
		return
	}
	fmt.Printf("%s: not found\n", val)

}

func FindInPath(bin string) (string, bool) {
	if file, exec := isExectutable(bin); exec {
		return file, true
	}
	paths := os.Getenv("PATH")
	arr := strings.Split(paths, ":")
	for _, path := range arr {
		fullpath := filepath.Join(path, bin)
		if file, err := isExectutable(fullpath); err {
			return file, true
		}
		if _, err := os.Stat(fullpath); err == nil {
			return fullpath, true
		}
	}
	return "", false
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

func ExtProg(argv []string) {
	if extCmd[argv[0]] {
		cmd := exec.Command(argv[0], argv[1:]...)
		// cmd.Args = argv.Args // Set argv to use original command name as argv[0]
		cmd.Args = argv
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {

			fmt.Fprintln(os.Stderr, "Error executing command:", err)
		}
		return
	}

	path, exists := isExectutable(argv[0])
	if exists || builtIns[path] {
		cmd := exec.Command(path, argv[1:]...)
		// cmd.Args = argv.Args // Set argv to use original command name as argv[0]
		cmd.Args = argv
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {

			fmt.Fprintln(os.Stderr, "Error executing command:", err)
		}

	} else {
		fmt.Printf("%s: command not found\n", argv[0])
	}

}

func isExectutable(filePath string) (string, bool) {

	// this will tell if the command exists in the path or not
	path, err := exec.LookPath(filePath)
	if err != nil {
		return "", false
	}
	return path, true
}

func Pwd() {
	dir, err := filepath.Abs(".")
	if err == nil {
		fmt.Println(dir)
	}
}

func Cd(argv []string) {
	if len(argv) < 2 {
		return
	}
	if argv[1] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "cd: could not find home directory")
			return
		}
		err = os.Chdir(homeDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", homeDir)
		}
		return
	}
	path := argv[1]
	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", path)
	}
}

func PrintArg(argv []string) {
	for _, arg := range argv {
		fmt.Println(arg)
	}

}

func SplitCmd(command string) []string {
	s := []string{}
	for _,ch := range command{
		if ch == '"'{
			s = charSplit(command,'"')
			return s
		}else if ch =='\''{
			s = charSplit(command,'\'')
			return s
		}
	}
	s = charSplit(command,'\'')
	return s
}

func charSplit(command string, ch byte) []string {

	s := []string{}
	flag := false
	curr := ""
	n := len(command)
	for i := 0; i < n-1; i++ {

		if command[i] == ' ' && !flag {
		
			if curr != "" {
				s = append(s, curr)
				curr = ""
			}
		
		} else if command[i] == ch {
		
			flag = !flag
		
		} else {
			curr += string(command[i])
		}

	}
		
		if curr != "" {
		s = append(s, curr)
		}
		
		return s

	
}	


