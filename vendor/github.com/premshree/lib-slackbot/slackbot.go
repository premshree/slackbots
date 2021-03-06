package slackbot

import(
  "log"
  "strings"

  "github.com/nlopes/slack"
)

const HELP = "help"

type Bot struct {
  api *slack.Client
  commands map[string]command
}

type command struct {
  Name string
  Description string
  Callback fn
}

type fn func(*Bot, string, string, ...string)

var (
  channelsMap map[string]interface{}
  usersMap map[string]string
)

// Initializes a new slackbot
func New(slackToken string) *Bot {
  return &Bot{
    api: slack.New(slackToken),
    commands: map[string]command{ },
  }
}

// AddCommand lets you add a command that your slack bot can respond to. It passes back
// the bot (*slackbot.Bot), a channel ID (string), a channel (string).
func (b *Bot) AddCommand(message, description string, callback fn) {
  b.commands[message] = command{
    Name: message,
    Description: description,
    Callback: callback,
  };
}

// Once you add commands to your bot, you need to call Run() so your bot can start
// listening to commands
func (b *Bot) Run() {
  rtm := b.api.NewRTM()
  go rtm.ManageConnection()

  channelsMap = b.getAllChannels()
  usersMap = b.getAllUsers()

  for msg := range rtm.IncomingEvents {
    switch ev := msg.Data.(type) {
    case *slack.MessageEvent:
      go b.handleMessage(ev.Msg)
    case *slack.RTMError:
      log.Printf("Error: %s\n", ev.Error())
    default:
    }
  }
}

// A handy function you can use within your AddCommand callbacks so the bot
// can reply to commands
func (b *Bot) Reply(channel string, reply string) {
  _, _, err := b.api.PostMessage(channel, reply, slack.PostMessageParameters{})
  if err != nil {
    log.Fatal(err)
  }
}

func (b *Bot) handleMessage(msg slack.Msg) {
  messageSlice := strings.Split(msg.Text, " ")
  command := messageSlice[0]
  channelID := msg.Channel
  var channelName string
  switch v := channelsMap[channelID].(type) {
  case slack.Channel:
    channelName = v.Name
  case slack.Group:
    channelName = v.Name
  }
  var args []string
  if len(messageSlice) > 1 {
    args = messageSlice[1:]
  }
  if _, ok := b.commands[command]; ok {
    log.Printf("♔ %s on #%s by @%s", command, channelName, b.Users()[msg.User])
    if args != nil && args[0] == HELP {
      b.Reply(channelID, b.commands[command].Description)
    } else {
      b.commands[command].Callback(b, channelID, channelName, args...)
    }
  }
}

func (b *Bot) Users() map[string]string {
  return usersMap
}

func (b *Bot) getAllChannels() map[string]interface{} {
  allChannels, err := b.api.GetChannels(true)
  if err != nil {
    log.Fatalf("Uh oh, error fetching channels: %v", err)
  }
  allGroups, err := b.api.GetGroups(true)
  if err != nil {
    log.Fatalf("Uh oh, error fetching private channels %v", err)
  }
  channelsMap := make(map[string]interface{})
  for _, channel := range allChannels {
    channelsMap[channel.ID] = channel
  }
  for _, group := range allGroups {
    channelsMap[group.ID] = group
  }

  return channelsMap
}

func (b *Bot) getAllUsers() map[string]string {
  allUsers, err := b.api.GetUsers()
  if err != nil {
    log.Fatalf("Uh oh, error fetching users: %v", err)
  }
  usersMap := make(map[string]string)
  for _, user := range allUsers {
    usersMap[user.ID] = user.Name
  }

  return usersMap
}
