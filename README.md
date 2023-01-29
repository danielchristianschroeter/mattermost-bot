# mattermost-bot

### What is it ?
mattermost-bot can join any Mattermost channel and watches for commands to execute them locally on the same host or call other remote API services. For security reasons, not all commands of all roles are available to use.

Use `!help` to list all possible commands (depends on configured ROLES variable) you can use in the Mattermost channel:

```
The following commands are availble to use:
[OpenShift Client]
!oc get     Display one or many resources
!oc describe     Show details of a specific resource or group of resources
!oc delete pod     Delete a pod
!oc scale     Set a new size for a deployment, replica set, or replication controller
!oc logs     Print the logs for a container in a pod
!oc adm top     Show usage statistics of resources on the server
!oc adm upgrade     Check on upgrade status or upgrade the cluster to a newer version
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
-help     Show help for a specified subcommand.
-version     An alias for the version subcommand.
[Terraform Wrapper]
!tf plan <ENV> <ZONE>     Show changes required by the current configuration for a zone within an environment.
!tf apply <ENV> <ZONE>     Create or update infrastructure for a zone within an environment.
[Govc Client]
!govc about     Display About info
!govc cluster.usage     Cluster resource usage summary
!govc datacenter.info     Information about datacenter
!govc datastore.cluster.info     Display datastore cluster info
!govc device.info     Device info for VM
!govc events     Display events
!govc find     Find managed objects
!govc metric.sample     Display metrics of host
!govc vm.info     Display info for VM
!govc vm.power     VM power operations
Global options:
-help     Show help for a specified subcommand.
```

If the output exceeds the message limit from Mattermost, the full output will automatically uploaded to a custom PrivateBin server (https://privatebin.info/).

### Configuration

`MB_DEBUG`: Show more details, for better debugging puposes you should enable this value.

`MB_ROLES`: Define the roles you want to use, seperated by comma (Available roles are: govc, oc, terraform, tf)

`MB_MATTERMOST_URL`: The Mattermost URL with protocol.

`MB_MATTERMOST_USERTOKEN`: The Mattermost personal access user token of a user the application can use to watch for commands in a channel.

`MB_MATTERMOST_TEAM`: The Mattermost team the user account is assigned to.

`MB_MATTERMOST_CHANNEL`: The Mattermost channel where the bot should watch and response for commands.

`MB_PRIVATEBIN_ENABLE`: Enable or disable (true/false) Privatebin (https://github.com/PrivateBin/PrivateBin), if the response exceeds 6315 characters. 

`MB_PRIVATEBIN_HOST`: The host of the Privatebin instance.

`MB_PRIVATEBIN_FORMATTER`: You can define the Privatebin formatter. (Recommendation: syntaxhighlighting).

`MB_PRIVATEBIN_EXPIRE`: Set the expire date when the links should be removed from Privatebin. (Recommendation: 6days)

`MB_PRIVATEBIN_OPENDISCUSSION`: Enable or disable (true/false) opendiscussion for Privatebin link. (Recommendation: false)

`MB_PRIVATEBIN_BURNAFTERREADING`: Enable or disable (true/false) burnafterreading for Privatebin link. (Recommendation: false)

`MB_PRIVATEBIN_PASSWORD`: Set a Privatebin password for the generated link. (Recommendation: leave the value empty)

`MB_GOVC_HOST`: Host from your vCenter Server.

`MB_GOVC_DATACENTER`: Datacenter of your vCenter Server.

`MB_GOVC_USERNAME`: Username with proper permissions to control vCenter Server.

`MB_GOVC_PASSWORD`: Password from the vCenter Server user account definied above.

`MB_TERRAFORM_EXECUTABLE`: Location of the terraform executable (https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)

`MB_MB_TERRAFORM_WRAPPER_EXECUTABLE`: Location of the terraform wrapper script executable (misc/tf.sh)

`MB_OC_EXECUTABLE`: Location of the OpenShift client executable (From your OpenShift Dashboard or https://github.com/okd-project/okd)

`MB_GOVC_EXECUTABLE`: Location of the govc executable (https://github.com/vmware/govmomi)

The configuration is located in the `config.env` file (you can rename the sample file `config.env.sample`).

The application will automatically use the environment variables if present, otherwise it will use the values from the config.env file.

### How you can build it?

You can build this Go project to a binary file on Ubuntu as follows:

* Install go from https://go.dev/dl/
* git clone https://github.com/danielchristianschroeter/mattermost-bot
* cd mattermost-bot
* go build .

You can also use the pre build binary from the releases page.

### How to create a proper Mattermost account?

Just create a normal user account for the bot (Bot Accounts are not working to read messages from a channel).
In the Profile page from this account, you need to create a "Personal Access Token". If you can not see this point in the Security menu, you have to enable this Option first in the System Console > Integration Management > Set "Enable Personal Access Tokens" to true.

### How can you run it ?

Rename `config.env.sample` to `config.env`, update the configuration and put the config file in the same directory as the executable.

Execute the mattermost-bot binary and see if the bot joining the defined Mattermost channel and respond to `!help`.

If you want to put the bot in the background, you can use the systemd example in `misc/mattermost-bot.service`.

Depending on your needs, the mattermost bot maybe nrequires higher permissions to execute certain commands.

I tested the Mattermost-bot ony with Ubuntu OS, but it whould work on other OS too.

### Avilable roles

At the moment, you can use the following roles:

terraform
tf (custom bash wrapper script for Terraform `misc/tf.sh`)
oc
govc

Just add or remove the role the bot should mange in the `ROLES=` variable in your configuration file or environment parameter.

Please note, the bot can only manage the role if the binary file and a proper configuration exist for that application.
