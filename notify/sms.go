package notify

type SmsSettings struct {
	Sms string `json:"sms"`
}

type SmsNotifier struct {
	Settings *SmsSettings
}

func (s *SmsNotifier) Notify(text string) error {
	return nil
}

func (ss *SmsSettings) Validate() error {
	return nil
}

func (s *SmsNotifier) String() string {
	return ""
}
