package oc

import (
	"errors"
	"fmt"
	"mattermost-bot/confighandler"
	"mattermost-bot/utils"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/exp/slices"
)

func ValidCommand(words []string, message string) (string, bool) {
	var reason string
	var valid bool
	validcommands := []string{"get", "describe", "delete", "scale", "logs", "adm", "api-versions", "api-resources", "help", "version", "options"}
	// command words lengths must be greater than 2
	if words[0] == "!oc" && len(words) >= 2 {
		if !utils.StringInSlice(words[1], validcommands) {
			reason = words[1] + " is not allowed with " + words[0]
			valid = false
			return reason, valid
		}
		// Match TRUSTED words (get, scale ...), prevent -f and -it
		if words[1] == "logs" && utils.StringInSlice("-f", words) {
			reason = "-f is not allowed with " + words[0] + " " + words[1]
			valid = false
			return reason, valid
		}
		if words[1] == "exec" && utils.StringInSlice("-it", words) {
			reason = "-it is not allowed with " + words[0] + " " + words[1]
			valid = false
			return reason, valid
		}
		if words[1] == "adm" && !(slices.Contains(words, "top") || slices.Contains(words, "upgrade")) {
			reason = "only top and upgrade is allowed with " + words[0] + " " + words[1]
			valid = false
			return reason, valid
		}
		if words[1] == "delete" && !slices.Contains(words, "pod") {
			reason = "only pod is allowed with " + words[0] + " " + words[1]
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
	executable := confighandler.App.Config.MB_OC_EXECUTABLE
	// Check if executable exist
	_, error := os.Stat(executable)
	if os.IsNotExist(error) {
		//fmt.Printf("%v does not exist.\n", executable)
		confighandler.App.Logger.Info().Str("function", "oc_Execute").Str("type", "response").Msg(executable + " does not exist.")
		return executable + " does not exist.", nil
	}
	cmd := strings.Replace(message, "!oc", executable, -1)
	// Check if command is a valid command
	reason, valid := ValidCommand(words, message)
	if valid {
		confighandler.App.Logger.Info().Str("function", "oc_Execute").Str("type", "request").Msg(cmd)
		args := strings.Split(cmd, " ")
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			confighandler.App.Logger.Error().Err(err).Str("function", "oc_Execute").Str("type", "response").Msg(string(output))
			cmdout = fmt.Sprintf("%s \n %s", err, output)
		} else {
			confighandler.App.Logger.Info().Str("function", "oc_Execute").Str("type", "response").Msg(string(output))
			cmdout = string(output)
		}
		return cmdout, nil
	} else {
		return reason, errors.New("Error in response: " + reason)
	}
}

func Help() string {
	available_commands := "[OpenShift Client]\n"
	available_commands += "!oc get\t Display one or many resources\n"
	available_commands += "!oc describe\t Show details of a specific resource or group of resources\n"
	available_commands += "!oc delete pod\t Delete a pod\n"
	available_commands += "!oc scale\t Set a new size for a deployment, replica set, or replication controller\n"
	available_commands += "!oc logs\t Print the logs for a container in a pod\n"
	available_commands += "!oc adm top\t Show usage statistics of resources on the server\n"
	available_commands += "!oc adm upgrade\t Check on upgrade status or upgrade the cluster to a newer version\n"
	available_commands += "!oc api-versions\t Print the supported API versions on the server, in the form of group/version\n"
	available_commands += "!oc api-resources\t Print the supported API resources on the server\n"
	available_commands += "!oc help\t Help about any command\n"
	available_commands += "!oc version\t Print the client and server version information\n"
	available_commands += "!oc options\t List of global command-line options (applies to all commands)\n"
	return available_commands
}
