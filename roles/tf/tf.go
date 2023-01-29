package tf

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
	KEYCOMMAND = "!tf"
)

func ValidCommand(words []string, message string) (string, bool) {
	var reason string
	var valid bool
	validcommands := []string{"plan", "apply"}
	// Command words lengths must be equal or greater than 2
	if words[0] == KEYCOMMAND && len(words) >= 2 {
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

func Execute(words []string, message string) (string, error) {
	var cmdout string
	executable := confighandler.App.Config.MB_TERRAFORM_WRAPPER_EXECUTABLE
	// Check if executable exist
	_, error := os.Stat(executable)
	if os.IsNotExist(error) {
		confighandler.App.Logger.Info().Str("function", "tf_Execute").Str("type", "response").Msg(executable + " does not exist.")
		return executable + " does not exist.", nil
	}
	cmd := strings.Replace(message, KEYCOMMAND, executable, -1)
	// Check if command is a valid command
	reason, valid := ValidCommand(words, message)
	if valid {
		confighandler.App.Logger.Info().Str("function", "tf_Execute").Str("type", "request").Msg(cmd)
		args := strings.Split(cmd, " ")
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			confighandler.App.Logger.Error().Err(err).Str("function", "tf_Execute").Str("type", "response").Msg(string(output))
			cmdout = fmt.Sprintf("%s \n %s", err, output)
		} else {
			confighandler.App.Logger.Info().Str("function", "tf_Execute").Str("type", "response").Msg(string(output))
			cmdout = string(output)
		}
		return cmdout, nil
	} else {
		return reason, errors.New("Error in response: " + reason)
	}
}

func Help() string {
	available_commands := "[Terraform Wrapper]\n"
	available_commands += KEYCOMMAND + " plan <ENV> <ZONE>\t Show changes required by the current configuration for a zone within an environment.\n"
	available_commands += KEYCOMMAND + " apply <ENV> <ZONE>\t Create or update infrastructure for a zone within an environment.\n"
	return available_commands
}
