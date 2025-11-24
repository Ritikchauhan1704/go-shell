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

func TypFun(argv []string, fileWriting bool, desc int, mode rune) {

	if len(argv) == 1 {
		return
	}

	val := argv[1]
	outputString := ""
	if builtIns[val] {
		outputString = fmt.Sprintf("%s is a shell builtin\n", val)

	} else if file, exists := findInPath(val); exists {
		outputString = fmt.Sprintf("%s is %s\n", val, file)

	}
	outputString = fmt.Sprintf("%s: not found\n", val)
	if fileWriting {
		if desc == 1 {

			writeToFile(argv[len(argv)-1], outputString, mode)
			return
		} else {
			writeToFile(argv[len(argv)-1], "", mode)
		}
	}
	fmt.Printf(outputString)

}

func ExitCmd(argv []string, fileWriting bool, desc int, mode rune) {
	code := 0
	if len(argv) > 1 {
		val, err := strconv.Atoi(argv[1])
		if err != nil {
			code = val
		} else if fileWriting && desc == 2 {
			writeToFile(argv[len(argv)-1], err.Error(), mode)
		}
	}
	os.Exit(code)
}

func EchoCmd(argv []string, fileWriting bool, desc int, mode rune) {
	output := strings.Join(argv[1:len(argv)-1], " ")
	if fileWriting {
		if desc == 1 {
			writeToFile(argv[len(argv)-1], output, mode)
			return
		} else {
			writeToFile(argv[len(argv)-1], "", mode)

		}

	}
	fmt.Println(output)
}

func ExtProg(argv []string, fileWriting bool, desc int, mode rune) {
	if extCmd[argv[0]] {
		cmd := exec.Command(argv[0], argv[1:len(argv)-1]...)
		// cmd.Args = argv.Args // Set argv to use original command name as argv[0]
		cmd.Args = argv
		cmd.Stdin = os.Stdin
		if fileWriting {
			file,_:= createFile(argv[len(argv)-1],mode)
			
			if  desc==1{
				cmd.Stdout = file
				cmd.Stderr = os.Stderr
			}else {
				cmd.Stderr = file
				cmd.Stdout = os.Stdout
			}
			// cmd.Stdout
		}
		if err := cmd.Run(); err != nil {
			errorOp := "Error executing command:" + err.Error()
			if fileWriting && desc == 2 {
				writeToFile(argv[len(argv)-1], errorOp, mode)
				return
			}
			fmt.Fprintln(os.Stderr, errorOp)
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

			errorOp := "Error executing command:" + err.Error()
			if fileWriting && desc == 2 {
				writeToFile(argv[len(argv)-1], errorOp, mode)
				return
			}
			fmt.Fprintln(os.Stderr, errorOp)
		}

	} else {
		output := fmt.Sprintf("%s: command not found\n", argv[0])
		if fileWriting && desc==1{
			writeToFile(argv[len(argv)-1],output,mode)
		}else{
		fmt.Printf(output)
		}

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

func Pwd(filename string, fileWriting bool, desc int, mode rune) {
	dir, err := filepath.Abs(".")

	if err == nil {
		fmt.Println(dir)
		if fileWriting && desc == 1 {
			writeToFile(filename, dir, mode)
		}
	} else if fileWriting && desc == 2 {
		writeToFile(filename, err.Error(), mode)
	}

}

func Cd(argv []string, fileWriting bool, desc int, mode rune) {
	if len(argv) < 2 {
		return
	}
	if argv[1] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			if fileWriting && desc==1{
				writeToFile(argv[len(argv)-1],"cd: could not find home directory",mode)
				return
			}
			fmt.Fprintln(os.Stderr, "cd: could not find home directory")
			return
		}
		err = os.Chdir(homeDir)
		if err != nil {
			op := fmt.Sprintf( "cd: %s: No such file or directory\n", homeDir)
			if fileWriting && desc==1{
				writeToFile(argv[len(argv)-1],op,mode)
			}else{
			fmt.Fprintf(os.Stderr,op)
			}
		}
		return
	}
	path := argv[1]
	err := os.Chdir(path)
	if err != nil {
		op := fmt.Sprintf( "cd: %s: No such file or directory\n", path)
			if fileWriting && desc==1{
				writeToFile(argv[len(argv)-1],op,mode)
			}else{
			fmt.Fprintf(os.Stderr,op)
			}
	}
}


var escCh = map[byte]bool{'"': true, '\\': true, '$': true, '`': true}

func SplitCmd(command string) ([]string, bool, int, rune) {

	s := []string{}
	singleQ, doubleQ, esc := false, false, false
	curr := ""
	writeFileflag := false
	mode := 'o'
	desc := 1
	n := len(command)
	for i := 0; i < n-1; i++ {
		ch := command[i]
		if !singleQ && !doubleQ && ch == '>' {
			writeFileflag = true
			if writeFileflag {
				if mode == 'a' {
					fmt.Errorf("error near '>'")
					return []string{}, false, -1, 'x'
				}
				mode = 'a'
			}
			if i > 0 && command[i-1] == '2' {
				desc = 2

			}
		}
		if esc && doubleQ {
			if !escCh[ch] {
				curr += "\\"

			}
			curr += string(ch)
			esc = false
		} else if esc {
			esc = false
			curr += string(ch)
		} else if ch == '\'' && !doubleQ {
			singleQ = !singleQ
		} else if ch == '"' && !singleQ {
			doubleQ = !doubleQ
		} else if ch == '\\' && !singleQ {
			esc = true
		} else if ch == ' ' && !singleQ && !doubleQ {
			if curr != "" {
				s = append(s, curr)
				curr = ""
			}
		} else {
			curr += (string)(ch)
		}

	}

	if curr != "" {
		s = append(s, curr)
	}

	return s, writeFileflag, desc, mode

}

func findInPath(bin string) (string, bool) {
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

// file creation based on mode
func createFile(filepath string, mode rune) (*os.File, error) {
	var file *os.File
	var err error
	if mode == 'a' {
		file, err = os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(filepath)
	}
	return file, err

}

// write to file function
func writeToFile(filename string, outputString string, mode rune) {
	file, err := createFile(filename, mode)
	if err == nil {
		file.WriteString(outputString)
		defer file.Close()
	}
}
