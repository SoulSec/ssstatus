package notify

import (
	"fmt"

	"github.com/gregdel/pushover"
)

type PushoverSettings struct {
	UserKey  string `json:"userKey"`
	AppToken string `json:"appToken"`
}

type PushoverNotifier struct {
	Settings *PushoverSettings
}

func (s *PushoverNotifier) Notify(text string) error {
	app := pushover.New(s.Settings.AppToken)
	recipient := pushover.NewRecipient(s.Settings.UserKey)
	message := pushover.NewMessageWithTitle(text+"not reached", "SSSTATUS Notification")
	_, err := app.SendMessage(message, recipient)
	return err
}

func (e *PushoverNotifier) Initialize() {

}
func (ts *PushoverSettings) Validate() error {
	errPushoverProperty := func(property string) error {
		return fmt.Errorf("Missing Pushover property %s", property)
	}
	switch {
	case ts.UserKey == "":
		return errPushoverProperty("user_key")
	case ts.AppToken == "":
		return errPushoverProperty("app_token")
	}
	return nil
}

func (t *PushoverNotifier) String() string {
	return fmt.Sprintf("Pushover: %s with appToken %s", t.Settings.UserKey, t.Settings.AppToken)
}
