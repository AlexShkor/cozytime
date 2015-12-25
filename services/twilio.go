package twilio

import (
  "fmt"
   "strconv"
  twilio "github.com/carlosdp/twiliogo"
  "bitbucket.org/AlexShkor/cozytime/settings"
  "math/rand"
)

func SendCode() (string, error) {
  conf := settings.Get()
  client := twilio.NewClient(conf.TwilioSID, "9dce779fef4e4988a3d721f401019d51")

  code := strconv.Itoa(rand.Intn(999999))
  fmt.Println("Code generated: " + code)
  message, err := twilio.NewMessage(client, "+15005550006", "+375259005003", twilio.Body("Your code: " + code))

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(message.Status)
    fmt.Println(message.Sid)
  }
  return code, err;
}