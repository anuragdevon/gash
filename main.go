package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Main function => driver code
func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		// Read the input from the user
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the command
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// Function to esxecute the command
func execInput(input string) error {

	// Remove the newline character
	input = strings.TrimSuffix(input, "\n")

	// Prepare the command to execute
	cmd := exec.Command(input)

	// Set the correct output device
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and return the error
	return cmd.Run()
}

// Function to esxecute the command
// func executeCommand(command string) {

// 	// Create a new command
// 	cmd := exec.Command("cmd", "/c", command)

// 	// Create a new output buffer
// 	var out bytes.Buffer

// 	// Set the output buffer to the command
// 	cmd.Stdout = &out

// 	// Run the command
// 	err := cmd.Run()
// 	if err ~= nil {
// 		fmt.Fprint(os.Stderr, "There was an error running the command: %s\n", err)
// 	}

// 	// Print the output
// 	fmt.Printf("%s\n", out.String())
// }
