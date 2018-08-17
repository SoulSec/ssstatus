package notify

type Initializer interface {
	Initialize()
}
type Notifier interface {
	Notify(text string) error
}

type Notifiers []Notifier

func (notifiers Notifiers) NotifyAll(text string) {
	for _, notifier := range notifiers {
		go notifier.Notify(text)
	}
}
