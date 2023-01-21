package terraform

import (
	"errors"
	"fmt"
	"mattermost-bot/utils"
	"os"
	"os/exec"
	"strings"
)

const (
	KEYCOMMAND = "!terraform"
)

func ValidCommand(words []string, message string) (string, bool) {
	var reason string
	var valid bool
	validcommands := []string{"init", "validate", "plan", "apply", "output", "show", "version"}
	// Command words lengths must be equal or greater than 2
	if words[0] == KEYCOMMAND && len(words) >= 2 {
		if !utils.StringInSlice(words[1], validcommands) && words[1][0:strings.Index(words[1], "=")] != "-chdir" {
			reason = words[1] + " is not allowed with " + words[0]
			valid = false
			return reason, valid
		}
		valid = true
	} else {
		reason = "length of words is " + fmt.Sprint(len(words))
		valid = false
	}
	return reason, valid
}

// TODO Add --auto-approve on apply
func Execute(words []string, message string) (string, error) {
	var cmdout string
	executable := "/usr/bin/terraform"
	// Check if executable exist
	_, error := os.Stat(executable)
	if os.IsNotExist(error) {
		fmt.Printf("%v does not exist.\n", executable)
		return executable + " does not exist.", nil
	}
	cmd := strings.Replace(message, KEYCOMMAND, executable, -1)
	cmd += " -no-color"
	// Check if command is a valid command
	reason, valid := ValidCommand(words, message)
	if valid {
		fmt.Printf("Request: %v\n", cmd)
		args := strings.Split(cmd, " ")
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Print stderr on error
			fmt.Println("Response:")
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			cmdout = fmt.Sprintf("%s \n %s", err, output)
		} else {
			fmt.Println("Response:")
			fmt.Println(string(output))
			cmdout = string(output)
		}
		return cmdout, nil
	} else {
		return reason, errors.New("Error in response: " + reason)
	}
}

func Help() string {
	available_commands := "[Terraform]\n"
	available_commands += KEYCOMMAND + " init\t Prepare your working directory for other commands\n"
	available_commands += KEYCOMMAND + " validate\t Check whether the configuration is valid\n"
	available_commands += KEYCOMMAND + " plan\t Show changes required by the current configuration\n"
	available_commands += KEYCOMMAND + " apply\t Create or update infrastructure\n"
	available_commands += KEYCOMMAND + " output\t Show output values from your root module\n"
	available_commands += KEYCOMMAND + " show\t Show the current state or a saved plan\n"
	available_commands += KEYCOMMAND + " version\t Show the current Terraform version\n"
	available_commands += "Global options:\n"
	available_commands += "-chdir=DIR\t Switch to a different working directory before executing the given subcommand.\n"
	available_commands += "-help\t Show help for a specified subcommand.\n"
	available_commands += "-version\t An alias for the version subcommand.\n"
	return available_commands
}
