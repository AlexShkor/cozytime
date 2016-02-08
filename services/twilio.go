package twilio

import (
  "fmt"
  twilio "github.com/carlosdp/twiliogo"
  "bitbucket.org/AlexShkor/cozytime/settings"
)

func SendCode(code string, phone string) (error) {
  conf := settings.Get()
  client := twilio.NewClient(conf.TwilioSID, "9dce779fef4e4988a3d721f401019d51")
  message, err := twilio.NewMessage(client, "+15005550006", "+" + phone, twilio.Body("Your code: " + code))

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(message.Status)
    fmt.Println(message.Sid)
  }
  return err;
}