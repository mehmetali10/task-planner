package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PromptForEnv asks for an environment variable with a default value
func PromptForEnv(key, defaultVal string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter value for %s (default: %s): ", key, defaultVal)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultVal
		}

		if IsValidEnvVar(key, input) {
			return input
		}

		fmt.Println("Invalid value. Please enter a valid input.")
	}
}

// IsValidEnvVar checks if the input is valid for a given environment variable
func IsValidEnvVar(key, value string) bool {
	switch key {
	case "DB_PORT":
		_, err := strconv.Atoi(value)
		return err == nil
	default:
		return len(value) > 0
	}
}

// PromptForInput asks for a general input from the user
func PromptForInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// PromptYesNo asks for a yes/no confirmation
func PromptYesNo(message string) bool {
	for {
		input := PromptForInput(message + " ")
		input = strings.ToLower(input)
		if input == "yes" || input == "y" || input == "" {
			return true
		} else if input == "no" || input == "n" {
			return false
		}
		fmt.Println("Please enter 'yes' or 'no'.")
	}
}
