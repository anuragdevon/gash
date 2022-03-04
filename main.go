package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var LINE_NUMBER int = 0
var HIST_FILE string = "gash_history.log"

const ClearLine = "\n\033[1A\033[K"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func promt() {
	path, err := os.Getwd()
	check(err)
	path = strings.Replace(string(path), "/home/anurag", "~", 1)
	colorGreen := "\033[32m"
	colorBlue := "\033[34m"
	colorReset := "\033[0m"
	fmt.Print(string(colorBlue), path, string(colorGreen), " > ", string(colorReset))
}

func unixSignals() {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT)

	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signalChanel
			fmt.Println("Signal Received: ", s)
			switch s {
			case syscall.SIGINT:
				fmt.Println("Signal interrupt triggered.")

			default:
				fmt.Println("Unknown signal.")
				exit_chan <- 1
			}
		}
	}()
	exitCode := <-exit_chan
	os.Exit(exitCode)
}

func editGashHistory(input string) {
	f, err := os.OpenFile(HIST_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	check(err)
	defer f.Close()
	_, errW := f.WriteString(input)
	check(errW)
}

func readGashHistory(lineNumber int) string {
	f, err := os.OpenFile(HIST_FILE, os.O_RDONLY, os.ModePerm)
	check(err)
	defer f.Close()

	rd := bufio.NewReader(f)
	i := 0
	for line, err := rd.ReadString('\n'); err != io.EOF; line, err = rd.ReadString('\n') {
		i += 1
		if lineNumber == i {
			return line
		}
	}
	return ""
}

func total_lines() int {
	f, err := os.OpenFile(HIST_FILE, os.O_RDONLY, os.ModePerm)
	check(err)
	defer f.Close()
	count := 0

	rd := bufio.NewReader(f)
	for _, err := rd.ReadString('\n'); err != io.EOF; _, err = rd.ReadString('\n') {
		count += 1
	}
	return count
}

func decisionTree(b []byte, executionStatus bool, prevCommand string) bool {
	if !executionStatus {
		var c []byte = make([]byte, 1)
		var d []byte = make([]byte, 1)

		if string(b) == string(byte(27)) {
			os.Stdin.Read(c)
			os.Stdin.Read(d)

			if string(c) == string(byte(91)) {
				if string(d) == string(byte(65)) {
					// read history
					LINE_NUMBER -= 1
					input := readGashHistory(LINE_NUMBER)
					input = strings.TrimSuffix(input, "\n")
					fmt.Print(ClearLine)
					promt()
					fmt.Print(input)
					prevCommand = input
					os.Stdin.Read(b)

					executionStatus = decisionTree(b, executionStatus, prevCommand)

				} else if string(d) == string(byte(66)) {
					// read latest
					LINE_NUMBER += 1
					input := readGashHistory(LINE_NUMBER)
					input = strings.TrimSuffix(input, "\n")
					fmt.Print(ClearLine)
					promt()
					fmt.Print(input)
					prevCommand = input
					os.Stdin.Read(b)

					executionStatus = decisionTree(b, executionStatus, prevCommand)
				}
			}
		} else {
			input := ""
			if prevCommand == "" {
				fmt.Print(string(b))

				// Enable chacter display on screen
				exec.Command("stty", "-F", "/dev/tty", "echo").Run()
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				check(err)
				input = string(b) + input

				editGashHistory(input)
				executionStatus = true

				if err = execInput(input); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

			} else {
				fmt.Print(string(b))
				input = prevCommand + string(b)
				// input = strings.TrimSuffix(input, "\n")
				// // Enable chacter display on screen
				// exec.Command("stty", "-F", "/dev/tty", "echo").Run()
				// reader := bufio.NewReader(os.Stdin)
				// extra_input, err := reader.ReadString('\n')
				// check(err)

				// input = input + extra_input

				editGashHistory(input)
				executionStatus = true

				if err := execInput(input); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}

		}
	}
	return executionStatus
}

func execInput(input string) error {

	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		// 'cd' to home dir with empty path not yet supported.
		if len(args) < 2 {
			dir := "/home/" + "anurag"
			return os.Chdir(dir)
		}
		// Change the directory and return the error
		return os.Chdir(args[1])

	case "exit":
		os.Exit(0)
	}
	cmd := exec.Command(args[0], args[1:]...)

	// Set the correct output device
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	var b []byte = make([]byte, 1)
	executionStatus := false
	for {
		// disble chacter display on screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		LINE_NUMBER = total_lines() + 1
		prevCommand := ""
		promt()
		os.Stdin.Read(b)
		decisionTree(b, executionStatus, prevCommand)
		// unixSignals()
	}
}
