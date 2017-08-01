package weather

import(
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"

  "github.com/premshree/lib-slackbot"
  "github.com/spf13/viper"
)

var (
  openWeatherMapToken string
)

func init() {
  viper := viper.New()
  viper.SetEnvPrefix("libslackbot")
  viper.AutomaticEnv()
  viper.ReadInConfig()
  openWeatherMapToken = viper.GetString("owm_token")
}

func Weather(bot *slackbot.Bot, channelID string, channelName string, args ...string) {
  if args == nil {
    bot.Reply(channelID, "Usage: ?weather zipcode")
    return
  }

  url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s", args[0], openWeatherMapToken)
  rs, err := http.Get(url)
  if err != nil {
    panic(err)
  }
  defer rs.Body.Close()

  bodyBytes, err := ioutil.ReadAll(rs.Body)
  if err != nil {
    panic(err)
  }

  var location, description string
  var temp, humidity float64
  c := make(map[string]interface{})
  err = json.Unmarshal(bodyBytes, &c)
  if err != nil {
    panic(err)
  }

  location = c["name"].(string)
  if val, ok := c["weather"].([]interface{}); ok {
    if val, ok := val[0].(map[string]interface{}); ok {
      description = val["description"].(string)
    }
  }
  if val, ok := c["main"].(map[string]interface{}); ok {
    temp = func(kelvin float64) float64 {
      return 1.8 * (kelvin - 273) +32
      }(val["temp"].(float64))
      humidity = val["humidity"].(float64)
    }

    bot.Reply(
      channelID,
      fmt.Sprintf("Weather in %s: %s, %d°F, %d%% humidity", location, description, int(temp), int(humidity)),
    )
  }
