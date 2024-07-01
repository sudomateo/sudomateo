package main

import (
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

type Notifier interface {
	Notify(message string) error
	Type() string
	Destination() string
}

// Don’t design with interfaces, discover them. - Rob Pike
// The bigger the interface, the weaker the abstraction. - Rob Pike
// Accept interfaces, return structs. - Jack Lindamood
func main() {
	notifiers := []Notifier{
		&DiscordClient{ChannelID: "#jobs"},
		&DiscordClient{ChannelID: "#meetups"},
		&EmailClient{Address: "sudomateo@example.com"},
		&EmailClient{Address: "matthew@example.com"},
	}

	SendNotification("New video uploaded!", notifiers)
}

func SendNotification(message string, notifiers []Notifier) {
	for _, n := range notifiers {
		logger := logger.With(
			"type", n.Type(),
			"destination", n.Destination(),
		)

		if err := n.Notify(message); err != nil {
			logger.Error("notification failed",
				"error", err,
			)

			// To process the next notification.
			continue
		}

		logger.Info("notification sent")
	}
}
