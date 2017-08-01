package slackbots

import(
  "bytes"
  "fmt"
  "log"
  "sort"

  "github.com/premshree/lib-slackbot"
  "github.com/PagerDuty/go-pagerduty"
  "github.com/spf13/viper"
)

const (
  CONFIG_ENV_PREFIX = "pagerduty_oncall"
  TPL_CHANNEL_NOT_CONFIGURED = "Uh oh, #%s is not configured for ?oncall"
)

type Config struct {
  Channels []ChannelConfig
}

type ChannelConfig struct {
  Name string `mapstructure:"name"`
  EscalationPolicyID string `mapstructure:"escalation_policy_id"`
}

var (
  config Config
  channelConfigMap map[string]ChannelConfig
  token string
)

func init() {
  viper := viper.New()
  viper.SetConfigFile("./config/pagerduty-oncall.json")
  viper.SetEnvPrefix(CONFIG_ENV_PREFIX)
  viper.AutomaticEnv()
  err := viper.ReadInConfig()
  if err != nil {
    log.Fatalf("Error reading config file: %v", err)
  }

  err = viper.Unmarshal(&config)
  if err != nil {
    log.Fatalf("unable to decode config into struct: %v", err)
  }

  token = viper.GetString("TOKEN")
  channelConfigMap = getChannelConfigMap()
}

func PagerDutyOnCall(bot *slackbot.Bot, channelID string, channelName string, args ...string) {
  var buffer bytes.Buffer
  var channelConfig ChannelConfig
  var ok bool
  if channelConfig, ok = channelConfigMap[channelName]; !ok {
    bot.Reply(channelID, fmt.Sprintf(TPL_CHANNEL_NOT_CONFIGURED, channelName))
    return
  }
  opts := pagerduty.ListOnCallOptions{
    EscalationPolicyIDs: []string{channelConfig.EscalationPolicyID},
  }

  client := pagerduty.NewClient(token)
  var escalationLevels []int
  if onCalls, err := client.ListOnCalls(opts); err != nil {
    log.Fatalf("Error listing on-calls for #%s", channelName)
  } else {
    escalationPolicyMap := getEscalationPolicyMap(onCalls.OnCalls)
    for e := range escalationPolicyMap {
      escalationLevels = append(escalationLevels, e)
    }
    sort.Ints(escalationLevels)
    for _, k := range escalationLevels {
      if k > 3 {
        break
      }
      buffer.WriteString(fmt.Sprintf("Level %d: %s\n", k, escalationPolicyMap[k]))
    }
    bot.Reply(channelID, buffer.String())
  }
}

func getEscalationPolicyMap(oncalls []pagerduty.OnCall) map[int]string {
  escalationPolicyMap := make(map[int]string, 0)
  for _, oncall := range oncalls {
    escalationPolicyMap[int(oncall.EscalationLevel)] = oncall.User.Summary
  }

  return escalationPolicyMap
}

func getChannelConfigMap() map[string]ChannelConfig {
  channelConfigMap := make(map[string]ChannelConfig)
  for _, channel := range config.Channels {
    channelConfigMap[channel.Name] = channel
  }

  return channelConfigMap
}
