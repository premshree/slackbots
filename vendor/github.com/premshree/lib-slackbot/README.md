# lib-slackbot

[![Build Status](https://travis-ci.org/premshree/lib-slackbot.svg?branch=master)](https://travis-ci.org/premshree/lib-slackbot)

lib-slackbot is a thin convenience wrapper around nlopes' excellent [Slack API in Go](https://github.com/nlopes/slack).

When writing [slack](https://slack.com/) bots in [Go](https://golang.org/) I found myself repeating a lot of the same boilerplate. This is where lib-slackbot is useful. Currently, lib-slackbot is useful whenever you want to write a bot that listens to commands with optional arguments like so:
```
?weather 11231
?oncall
```

See the `examples` directory for an implementation.

## Installation

```
$ go get github.com/premshree/lib-slackbot
```

## Example

```go
package main

import(
  "time"

  "github.com/premshree/lib-slackbot"
)

func main() {
  token := "YOUR_SLACK_BOT_TOKEN"
  bot := slackbot.New(token)

  bot.AddCommand("?time", "What's the local time?", func(bot *slackbot.Bot, channelID string, channelName string, args ...string) {
    t := time.Now()
    bot.Reply(channelID, t.Format("Mon Jan _2 15:04:05 2006"))
  })

  // Run starts the bot's listen loop
  bot.Run()
}

```
