package main

import (
	"fmt"
	"mattermost-bot/roles/govc"
	"mattermost-bot/roles/oc"
	"mattermost-bot/roles/terraform"
	"mattermost-bot/roles/tf"
	utils "mattermost-bot/utils"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

var version = "DEV"

var client *model.Client4
var webSocketClient *model.WebSocketClient

var botUser *model.User
var botTeam *model.Team
var botChannel *model.Channel

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {
	println("Mattermost-Bot (v" + version + ") started.")
	println("Loaded configuration:")
	fmt.Printf("MATTERMOST_HOST: %s\n", utils.GetConfigValue("MATTERMOST_HOST"))
	fmt.Printf("MATTERMOST_EMAIL: %s\n", utils.GetConfigValue("MATTERMOST_EMAIL"))
	fmt.Printf("MATTERMOST_USER: %s\n", utils.GetConfigValue("MATTERMOST_USER"))
	fmt.Printf("MATTERMOST_TEAM: %s\n", utils.GetConfigValue("MATTERMOST_TEAM"))
	fmt.Printf("MATTERMOST_CHANNEL: %s\n", utils.GetConfigValue("MATTERMOST_CHANNEL"))
	SetupGracefulShutdown()
	url := fmt.Sprintf("https://%s", utils.GetConfigValue("MATTERMOST_HOST"))
	client = model.NewAPIv4Client(url)

	// Lets test to see if the mattermost server is up and running
	MakeSureServerIsRunning()

	// lets attempt to login to the Mattermost server as the bot user
	// This will set the token required for all future calls
	// You can get this token with client.AuthToken
	LoginAsTheBotUser()

	// Get Team for Bot
	FindBotTeam()

	// Join Bot to Channel
	JoinBotChannel()

	// Lets start listening to some channels via the websocket!
	wssURL := fmt.Sprintf("wss://%s", utils.GetConfigValue("MATTERMOST_HOST"))
	for {
		webSocketClient, err := model.NewWebSocketClient4(wssURL, client.AuthToken)
		if err != nil {
			println("We failed to connect to the web socket")
			PrintError(err)
		}
		println("Successfully connected to websocket " + wssURL)
		webSocketClient.Listen()

		for resp := range webSocketClient.EventChannel {
			HandleWebSocketResponse(resp)
		}
	}
}

func MakeSureServerIsRunning() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		println("Mattermost Server detected and is running version " + props["Version"])
	}
}

func LoginAsTheBotUser() {
	if user, resp := client.Login(utils.GetConfigValue("MATTERMOST_EMAIL"), utils.GetConfigValue("MATTERMOST_PASSWORD")); resp.Error != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		botUser = user
	}
}

func FindBotTeam() {
	if team, resp := client.GetTeamByName(utils.GetConfigValue("MATTERMOST_TEAM"), ""); resp.Error != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + utils.GetConfigValue("MATTERMOST_TEAM") + "'")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		botTeam = team
	}
}

func JoinBotChannel() {
	// Get channel id
	if channel, resp := client.GetChannelByName(utils.GetConfigValue("MATTERMOST_CHANNEL"), botTeam.Id, ""); resp.Error != nil {
		println("Failed to get the channel by name. Is MATTERMOST_CHANNEL correct and exist?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		botChannel = channel
	}

	// Join channel
	if channel, resp := client.AddChannelMember(botChannel.Id, botUser.Id); resp.Error != nil {
		println("Failed to join the channel.")
		PrintError(resp.Error)
		fmt.Printf("%v", channel)
	} else {
		fmt.Printf("Joined #" + utils.GetConfigValue("MATTERMOST_CHANNEL") + "\n")
	}
}

func SendMsgToChannel(msg string, replyToId string, codeblock bool) {
	// Limit output to 6315 characters
	shortened_msg, err := utils.CheckCharLimit(msg, 6315, codeblock)
	if err != nil {
		println(err)
	}
	post := model.Post{
		ChannelId: botChannel.Id,
		Message:   shortened_msg,
		RootId:    replyToId,
	}
	// TODO Remove clock, add check to msg https://api.mattermost.com/#operation/SaveReaction
	if _, resp := client.CreatePost(&post); resp.Error != nil {
		println("We failed to send a message to the logging channel")
		PrintError(resp.Error)
	}
}

func HandleWebSocketResponse(event *model.WebSocketEvent) {
	HandleMsgFromDebuggingChannel(event)
}

func HandleMsgFromDebuggingChannel(event *model.WebSocketEvent) {
	// If this isn't the debugging channel then lets ingore it
	if event.Broadcast.ChannelId != botChannel.Id {
		return
	}

	// Lets only reponded to messaged posted events
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {

		// ignore my events
		if post.UserId == botUser.Id {
			return
		}
		if matched, _ := regexp.MatchString("^!", post.Message); matched {
			message := strings.TrimSpace(post.Message)
			fmt.Printf("Command received: %v\n", message)
			words := strings.Fields(post.Message)
			// Fetch keycommand after exclamation mark
			keycommand := strings.Trim(words[0], "!")
			switch strings.Contains(utils.GetConfigValue("ROLES"), keycommand) {
			case true:
				// keycommand found in allowed roles
				switch keycommand {
				case "govc":
					SendMsgToChannel("Executing command...", post.Id, false)
					cmdout, err := govc.Execute(words, message)
					if err != nil {
						fmt.Println(err)
					}
					SendMsgToChannel(cmdout, post.Id, true)
				case "oc":
					SendMsgToChannel("Executing command...", post.Id, false)
					cmdout, err := oc.Execute(words, message)
					if err != nil {
						fmt.Println(err)
					}
					SendMsgToChannel(cmdout, post.Id, true)
				case "terraform":
					SendMsgToChannel("Executing command...", post.Id, false)
					cmdout, err := terraform.Execute(words, message)
					if err != nil {
						fmt.Println(err)
					}
					SendMsgToChannel(cmdout, post.Id, true)
				case "tf":
					SendMsgToChannel("Executing command...", post.Id, false)
					cmdout, err := tf.Execute(words, message)
					if err != nil {
						fmt.Println(err)
					}
					SendMsgToChannel(cmdout, post.Id, true)
				}
			case false:
				// No valid role found
				// if help command
				if strings.Contains("help", keycommand) {
					s := strings.Split(utils.GetConfigValue("ROLES"), ",")
					help := "The following commands are availble to use:\n"
					for _, s := range s {
						switch s {
						case "govc":
							help += govc.Help()
						case "oc":
							help += oc.Help()
						case "terraform":
							help += terraform.Help()
						case "tf":
							help += tf.Help()
						}
					}
					SendMsgToChannel(help, post.Id, false)
					return
				}
				// or unknown command
				SendMsgToChannel("Unknown command. Use !help to list all vailable commands.", post.Id, false)
			}
		}
		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)alive(?:$|\W)`, post.Message); matched {
			SendMsgToChannel("Yes, I'm running", post.Id, false)
			return
		}
	}
}

func PrintError(err *model.AppError) {
	println("\tError Details:")
	println("\t\t" + err.Message)
	println("\t\t" + err.Id)
	println("\t\t" + err.DetailedError)
}

func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if webSocketClient != nil {
				webSocketClient.Close()
			}
			// Leave channel
			if _, resp := client.RemoveUserFromChannel(botChannel.Id, botUser.Id); resp.Error != nil {
				println("Failed to remove user from channel.")
				PrintError(resp.Error)
			} else {
				fmt.Printf("Leave #" + utils.GetConfigValue("MATTERMOST_CHANNEL") + "\n")
			}

			os.Exit(0)
		}
	}()
}
