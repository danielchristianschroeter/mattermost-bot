package terraform

import (
	"errors"
	"fmt"
	"mattermost-bot/confighandler"
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
		// Check if -chdir is required and in the correct position (allow required for init, validate, plan or apply)
		if len(words) > 1 && strings.Contains(words[1], "-chdir=") {
			if !utils.StringInSlice(words[2], []string{"init", "validate", "plan", "apply"}) {
				reason = "-chdir only allowed for init, validate, plan or apply"
				valid = false
				return reason, valid
			} else {
				valid = true
				return reason, valid
			}
		}
		if !utils.StringInSlice(words[1], validcommands) {
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
	executable := confighandler.App.Config.MB_TERRAFORM_EXECUTABLE
	// Check if executable exist
	_, error := os.Stat(executable)
	if os.IsNotExist(error) {
		confighandler.App.Logger.Info().Str("function", "terraform_Execute").Str("type", "response").Msg(executable + " does not exist.")
		return executable + " does not exist.", nil
	}
	cmd := strings.Replace(message, KEYCOMMAND, executable, -1)
	cmd += " -no-color"
	// Check if command is a valid command
	reason, valid := ValidCommand(words, message)
	if valid {
		confighandler.App.Logger.Info().Str("function", "terraform_Execute").Str("type", "request").Msg(cmd)
		args := strings.Split(cmd, " ")
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			confighandler.App.Logger.Error().Err(err).Str("function", "terraform_Execute").Str("type", "response").Msg(string(output))
			cmdout = fmt.Sprintf("%s \n %s", err, output)
		} else {
			confighandler.App.Logger.Info().Str("function", "terraform_Execute").Str("type", "response").Msg(string(output))
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
