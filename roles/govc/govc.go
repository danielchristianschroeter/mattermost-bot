package govc

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
	KEYCOMMAND = "!govc"
)

func ValidCommand(words []string, message string) (string, bool) {
	var reason string
	var valid bool
	validcommands := []string{"about", "cluster.usage", "datacenter.info", "datastore.cluster.info", "device.info", "events", "find", "metric.sample", "vm.info", "vm.power"}
	// Command words lengths must be equal or greater than 2
	if words[0] == KEYCOMMAND && len(words) >= 2 {
		if !utils.StringInSlice(words[1], validcommands) {
			reason = words[1] + " is not allowed with " + words[0]
			valid = false
			return reason, valid
		}
		// Find without any additional commands is not allowed
		if words[1] == "find" && len(words) <= 2 {
			reason = "find without any filter not allowed. Use --help to show options."
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
	executable := confighandler.App.Config.MB_GOVC_EXECUTABLE
	// Check if executable exist
	_, error := os.Stat(executable)
	if os.IsNotExist(error) {
		confighandler.App.Logger.Info().Str("function", "govc_Execute").Str("type", "response").Msg(executable + " does not exist.")
		return executable + " does not exist.", nil
	}
	cmd := strings.Replace(message, KEYCOMMAND, executable, -1)
	// Check if command is a valid command
	reason, valid := ValidCommand(words, message)
	if valid {
		confighandler.App.Logger.Info().Str("function", "govc_Execute").Str("type", "request").Msg(cmd)
		os.Setenv("GOVC_URL", confighandler.App.Config.MB_GOVC_URL)
		os.Setenv("GOVC_DATACENTER", confighandler.App.Config.MB_GOVC_DATACENTER)
		os.Setenv("GOVC_USERNAME", confighandler.App.Config.MB_GOVC_USERNAME)
		os.Setenv("GOVC_PASSWORD", confighandler.App.Config.MB_GOVC_PASSWORD)
		args := strings.Split(cmd, " ")
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			confighandler.App.Logger.Error().Err(err).Str("function", "govc_Execute").Str("type", "response").Msg(string(output))
			cmdout = fmt.Sprintf("%s \n %s", err, output)
		} else {
			confighandler.App.Logger.Info().Str("function", "govc_Execute").Str("type", "response").Msg(string(output))
			cmdout = string(output)
		}
		return cmdout, nil
	} else {
		return reason, errors.New("Error in response: " + reason)
	}
}

func Help() string {
	available_commands := "[Govc Client]\n"
	available_commands += KEYCOMMAND + " about\t Display About info\n"
	available_commands += KEYCOMMAND + " cluster.usage\t Cluster resource usage summary\n"
	available_commands += KEYCOMMAND + " datacenter.info\t Information about datacenter\n"
	available_commands += KEYCOMMAND + " datastore.cluster.info\t Display datastore cluster info\n"
	available_commands += KEYCOMMAND + " device.info\t Device info for VM\n"
	available_commands += KEYCOMMAND + " events\t Display events\n"
	available_commands += KEYCOMMAND + " find\t Find managed objects\n"
	available_commands += KEYCOMMAND + " metric.sample\t Display metrics of host\n"
	available_commands += KEYCOMMAND + " vm.info\t Display info for VM\n"
	available_commands += KEYCOMMAND + " vm.power\t VM power operations\n"
	available_commands += "Global options:\n"
	available_commands += "-help\t Show help for a specified subcommand.\n"
	return available_commands
}
