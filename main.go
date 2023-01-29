package main

import (
	"encoding/json"
	"fmt"
	"math"
	"mattermost-bot/confighandler"
	"mattermost-bot/roles/govc"
	"mattermost-bot/roles/oc"
	"mattermost-bot/roles/terraform"
	"mattermost-bot/roles/tf"
	"mattermost-bot/utils"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
)

type Command struct {
	name      string
	help      string
	executeFn func(words []string, message string) (string, error)
}

var commands = map[string]Command{
	"govc": {
		name:      "govc",
		help:      govc.Help(),
		executeFn: govc.Execute,
	},
	"oc": {
		name:      "oc",
		help:      oc.Help(),
		executeFn: oc.Execute,
	},
	"terraform": {
		name:      "terraform",
		help:      terraform.Help(),
		executeFn: terraform.Execute,
	},
	"tf": {
		name:      "tf",
		help:      tf.Help(),
		executeFn: tf.Execute,
	},
}

func main() {
	confighandler.Init()

	setupGracefulShutdown()

	// Create a new mattermost client
	confighandler.App.MattermostClient = model.NewAPIv4Client(confighandler.App.Config.MB_MATTERMOST_URL)

	// Check if Mattermost server is reachable
	if props, resp, err := confighandler.App.MattermostClient.GetOldClientConfig(""); err != nil {
		confighandler.App.Logger.Fatal().Stack().Err(err).Str("function", "main").Interface("response", resp).Msg("There was a problem reaching the Mattermost server. Are you sure Mattermost url definied in MB_MATTERMOST_URL is correct and the Mattermost server is running?")
		os.Exit(1)
	} else {
		confighandler.App.Logger.Debug().Str("function", "main").Interface("response", resp).Interface("properties", props).Msg("Mattermost server detected")
		//println("Mattermost Server detected and is running version " + props["Version"])
	}

	// Login with user token
	confighandler.App.MattermostClient.SetToken(confighandler.App.Config.MB_MATTERMOST_USERTOKEN)

	if user, resp, err := confighandler.App.MattermostClient.GetUser("me", ""); err != nil {
		confighandler.App.Logger.Fatal().Stack().Err(err).Str("function", "main").Msg("Could not log in to Mattermost server. Are you sure the user token definied in MB_MATTERMOST_USERTOKEN is correct?")
	} else {
		confighandler.App.Logger.Debug().Str("function", "main").Interface("user", user).Interface("response", resp).Msg("")
		confighandler.App.Logger.Info().Str("function", "main").Msg("Successfully logged in to Mattermost server")
		confighandler.App.MattermostUser = user
	}

	// Find and save the bot's team to app struct
	if team, resp, err := confighandler.App.MattermostClient.GetTeamByName(confighandler.App.Config.MB_MATTERMOST_TEAM, ""); err != nil {
		confighandler.App.Logger.Fatal().Stack().Err(err).Str("function", "main").Interface("response", resp).Msg("Could not find team. Is this bot user account a member of the team? Are you sure the team definied in MB_MATTERMOST_TEAM is correct?")
	} else {
		confighandler.App.Logger.Debug().Str("function", "main").Interface("team", team).Interface("response", resp).Msg("")
		confighandler.App.MattermostTeam = team
	}

	// Find and save the talking channel to app struct
	if channel, resp, err := confighandler.App.MattermostClient.GetChannelByName(
		confighandler.App.Config.MB_MATTERMOST_CHANNEL, confighandler.App.MattermostTeam.Id, "",
	); err != nil {
		confighandler.App.Logger.Fatal().Stack().Err(err).Str("function", "main").Interface("response", resp).Msg("Could not find channel. Are you sure the channel defined in MB_MATTERMOST_CHANNEL exists?")
	} else {
		confighandler.App.Logger.Debug().Str("function", "main").Interface("channel", channel).Interface("response", resp).Msg("")
		confighandler.App.MattermostChannel = channel
	}

	// Join Bot to Channel
	JoinBotChannel()

	// Listen to live events coming in via websocket
	listenToEvents()
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if confighandler.App.MattermostWebSocketClient != nil {
				confighandler.App.Logger.Info().Str("function", "setupGracefulShutdown").Msg("Closing websocket connection to Mattermost server")
				confighandler.App.MattermostWebSocketClient.Close()
			}
			// Leave channel
			if _, err := confighandler.App.MattermostClient.RemoveUserFromChannel(confighandler.App.MattermostChannel.Id, confighandler.App.MattermostUser.Id); err != nil {
				confighandler.App.Logger.Fatal().Stack().Err(err).Str("function", "setupGracefulShutdown").Msg("Failed to remove bot user account from channel")
			} else {
				fmt.Printf("Leave #" + confighandler.App.Config.MB_MATTERMOST_CHANNEL + "\n")
			}

			confighandler.App.Logger.Info().Str("function", "setupGracefulShutdown").Msg("Shutting down")
			os.Exit(0)
		}
	}()
}

// func sendMsgToTalkingChannel(app *application, msg string, replyToId string) {
// 	// Note that replyToId should be empty for a new post.
// 	// All replies in a thread should reply to root.

// 	post := &model.Post{}
// 	post.ChannelId = app.mattermostChannel.Id
// 	post.Message = msg

// 	post.RootId = replyToId

// 	if _, _, err := app.mattermostClient.CreatePost(post); err != nil {
// 		confighandler.App.Logger.Error().Err(err).Str("RootID", replyToId).Msg("Failed to create post")
// 	}
// }

func sendMsgToChannel(msg string, replyToId string, codeblock bool) {
	// Limit output to 6315 characters
	shortened_msg, err := utils.CheckCharLimit(msg, 6315, codeblock)
	if err != nil {
		confighandler.App.Logger.Error().Stack().Err(err).Str("function", "sendMsgToChannel").Msg("CheckCharLimit failed for some reason")
	}
	post := &model.Post{}
	post.ChannelId = confighandler.App.MattermostChannel.Id
	post.Message = shortened_msg

	post.RootId = replyToId

	if _, _, err := confighandler.App.MattermostClient.CreatePost(post); err != nil {
		confighandler.App.Logger.Error().Stack().Err(err).Str("function", "sendMsgToChannel").Str("RootID", replyToId).Msg("Failed to send a message to the channel")
	}
}

func listenToEvents() {
	var err error
	failCount := 0
	wsURL := utils.CheckAndFormatURL(confighandler.App.Config.MB_MATTERMOST_URL)

	for {
		confighandler.App.MattermostWebSocketClient, err = model.NewWebSocketClient4(
			wsURL,
			confighandler.App.MattermostClient.AuthToken,
		)
		if err != nil {
			confighandler.App.Logger.Warn().Err(err).Str("function", "listenToEvents").Msg("Websocket connection to Mattermost server lost/disconnected, retrying...")
			failCount += 1
			// Wait for 2^failCount seconds before retrying to establish a connection to the Mattermost server through websocket
			time.Sleep(time.Duration(math.Pow(2, float64(failCount))) * time.Second)
			continue
		}
		confighandler.App.Logger.Info().Str("function", "listenToEvents").Msg("Websocket to Mattermost server established")

		confighandler.App.MattermostWebSocketClient.Listen()

		for event := range confighandler.App.MattermostWebSocketClient.EventChannel {
			// Launch new goroutine for handling the actual event
			// If required, you can limit the number of events beng processed at a time
			go handleWebSocketEvent(event)
		}
	}
}

func handleWebSocketEvent(event *model.WebSocketEvent) {

	// Ignore other channels
	if event.GetBroadcast().ChannelId != confighandler.App.MattermostChannel.Id {
		return
	}

	// Ignore other types of events
	if event.EventType() != model.WebsocketEventPosted {
		return
	}

	// Since this event is a post, unmarshal it to (*model.Post)
	post := &model.Post{}
	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
	if err != nil {
		confighandler.App.Logger.Error().Stack().Err(err).Str("function", "handleWebSocketEvent").Msg("Could not cast event to *model.Post")
	}

	// Ignore messages sent by this bot itself
	if post.UserId == confighandler.App.MattermostUser.Id {
		return
	}

	// Handle however you want
	handlePost(post)
}

func handlePost(post *model.Post) {
	confighandler.App.Logger.Debug().Str("function", "handlePost").Str("type", "request").Interface("post", post).Msg("")

	words := strings.Fields(post.Message)
	if len(words) == 0 || words[0][0] != '!' {
		return
	}
	keycommand := strings.Trim(words[0], "!")
	if command, ok := commands[keycommand]; ok {
		confighandler.App.Logger.Debug().Str("function", "handlePost").Str("type", "request").Interface("received_keycommand", keycommand)
		sendMsgToChannel("Executing command...", post.Id, false)
		cmdout, err := command.executeFn(words, post.Message)
		if err != nil {
			confighandler.App.Logger.Error().Stack().Err(err).Str("function", "handlePost").Msg("Executing command failed for some reason")
		}
		sendMsgToChannel(cmdout, post.Id, true)
	} else if keycommand == "help" {
		confighandler.App.Logger.Debug().Str("function", "handlePost").Str("type", "request").Interface("received_keycommand", keycommand)
		help := "The following commands are available to use:\n"
		for _, command := range commands {
			help += command.help
		}
		sendMsgToChannel(help, post.Id, false)
		return
	} else {
		sendMsgToChannel("Unknown command. Use !help to list all available commands.", post.Id, false)
	}

}

func JoinBotChannel() {
	// Join channel
	if channel, resp, err := confighandler.App.MattermostClient.AddChannelMember(confighandler.App.MattermostChannel.Id, confighandler.App.MattermostUser.Id); err != nil {
		confighandler.App.Logger.Error().Stack().Err(err).Str("function", "JoinBotChannel").Msg("Failed to join the channel")
		confighandler.App.Logger.Debug().Str("function", "JoinBotChannel").Interface("channel", channel).Interface("resp", resp).Msg("")
	} else {
		confighandler.App.Logger.Info().Str("function", "JoinBotChannel").Msg("Joined #" + confighandler.App.Config.MB_MATTERMOST_CHANNEL)
	}
}
