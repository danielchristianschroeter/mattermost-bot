# mattermost-bot

### What is it ?
mattermost-bot can join any Mattermost channel and watches for commands to execute them locally on the same host or call other remote API services. For security reasons, not all commands of all roles are available to use.

Use `!help` to list all possible commands (depends on configured ROLES variable) you can use in the channel:

```
The following commands are availble to use:
[OpenShift Client]
!oc get     Display one or many resources
!oc describe     Show details of a specific resource or group of resources
!oc scale     Set a new size for a deployment, replica set, or replication controller
!oc logs     Print the logs for a container in a pod
!oc adm top     Show usage statistics of resources on the server
!oc api-versions     Print the supported API versions on the server, in the form of group/version
!oc api-resources     Print the supported API resources on the server
!oc help     Help about any command
!oc version     Print the client and server version information
!oc options     List of global command-line options (applies to all commands)
[Terraform]
!terraform init     Prepare your working directory for other commands
!terraform validate     Check whether the configuration is valid
!terraform plan     Show changes required by the current configuration
!terraform apply     Create or update infrastructure
!terraform output     Show output values from your root module
!terraform show     Show the current state or a saved plan
!terraform version     Show the current Terraform version
Global options:
-chdir=DIR     Switch to a different working directory before executing the given subcommand.
-help     Show this help output, or the help for a specified subcommand.
-version     An alias for the version subcommand.
```

If the output exceeds the message limit from Mattermost, the full output will automatically uploaded to a custom PrivateBin server (https://privatebin.info/).

The configuration is located in the `config.env` file (you can rename the sample file `config.env.sample`).

All parameters can also be set as environment variables.

### How you can build it?

You can build this Go project to a binary file on Ubuntu as follows:

* Install go from https://go.dev/dl/
* git clone https://github.com/danielchristianschroeter/mattermost-bot
* cd mattermost-bot
* go build .

You can also use the pre build binary from the releases page.

### How can you run it ?

Rename `config.env.sample` to `config.env`, update the configuration and put the config file in the same directory as the executable (Alternative location: /etc/mattermost-bot/config.env)

Execute the mattermost-bot binary and see if the bot joining the defined Mattermost channel and respond to `!help`.

If you want to put the bot in the background, you can use the systemd example in `misc/mattermost-bot.service`.

Please note that the bot requires higher permissions to execute certain commands.

Mattermost-bot is only tested on Ubuntu OS systems.

### Avilable roles

At the moment, you can use the following roles:

terraform
tf (custom bash wrapper script for Terraform `misc/tf.sh`)
oc
govc

Just add or remove the role the bot should mange in the `ROLES=` variable in your configuration file or environment parameter.

Please note, the bot can only manage the role if the binary file and a proper configuration exist for that application.
