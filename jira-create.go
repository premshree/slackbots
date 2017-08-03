package slackbots

import(
  "bytes"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "regexp"
  "strings"

  "github.com/premshree/lib-slackbot"
  "github.com/spf13/viper"
)

type JiraResponse struct {
  Key string `json:"key"`
}

const (
  JIRA_CREATE_PATTERN = "(^[\\w]+)[\\s]+([\\w\\s]+)[\\s]+@[\\s]+([A-Za-z_]+$)"
  USAGE = "?jiracreate YOURPROJECT summary @ asignee"
  JIRA_ENV_PREFIX = "JIRA_CREATE"
)

var (
  jiraAuth string // Uses basic auth: base64(username:password)
  jiraBaseUrl string
)

func init() {
  viper := viper.New()
  viper.SetEnvPrefix(JIRA_ENV_PREFIX)
  viper.AutomaticEnv()
  jiraAuth = viper.GetString("AUTH")
  jiraBaseUrl = viper.GetString("BASE_URL")
}

func JiraCreate(bot *slackbot.Bot, channelID string, channelName string, args ...string) {
  url := fmt.Sprintf("%s/rest/api/2/issue", jiraBaseUrl)
  if args == nil {
    bot.Reply(channelID, "Usage: ?jira KEY summary @asignee")
    return
  }

  argsString := strings.Join(args, " ")
  r := regexp.MustCompile(JIRA_CREATE_PATTERN)
  if !r.MatchString(argsString) {
    bot.Reply(channelID, fmt.Sprintf("Usage: %s", USAGE))
    return
  }
  matches := r.FindAllStringSubmatch(argsString, -1)
  var key, summary, asignee string
  key = strings.ToUpper(matches[0][1])
  summary = matches[0][2]
  asignee = matches[0][3]

  jsonTpl := `{
    "fields": {
      "project":
      {
        "key": "%s"
      },
      "summary": "%s",
      "description": "%s",
      "issuetype": {
        "name": "Bug"
      },
      "assignee": {
        "name": "%s"
      }
    }
  }`

  jsonStr := []byte(fmt.Sprintf(jsonTpl, key, summary, summary, asignee))
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", fmt.Sprintf("Basic %s", jiraAuth))

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    log.Fatalf("Error creating Jira ticket: %v", err)
  }
  defer resp.Body.Close()

  body, _ := ioutil.ReadAll(resp.Body)
  var ret JiraResponse
  if err := json.Unmarshal(body, &ret); err != nil {
    log.Printf("Error unmarshaling json: %f", err)
  }
  if ret.Key == "" {
    bot.Reply(channelID, "Usage: ?jira KEY summary @asignee")
    return
  }
  bot.Reply(channelID, fmt.Sprintf("Issue created: %s/browse/%s", jiraBaseUrl, ret.Key))
}
