package main

import(
  "github.com/premshree/lib-slackbot"
  "github.com/premshree/slackbots/pagerduty-oncall"
  "github.com/premshree/slackbots/weather"
  "github.com/spf13/viper"
)

var (
  slackToken string
)

func init() {
  viper := viper.New()
  viper.SetEnvPrefix("omnibot")
  viper.AutomaticEnv()
  viper.ReadInConfig()
  slackToken = viper.GetString("slack_token")
}

func main() {
  bot := slackbot.New(slackToken)

  bot.AddCommand("?oncall", "Who's on call", pagerduty.PagerDutyOnCall)
  bot.AddCommand("?weather", "Usage: ?weather zipcode", weather.Weather)

  bot.Run()
}
