package main

import "fmt"

type DiscordClient struct {
	ChannelID string
}

func (d *DiscordClient) Notify(message string) error {
	if err := simulateError(); err != nil {
		return fmt.Errorf("failed reaching discord: %w", err)
	}

	return nil
}

func (d *DiscordClient) Type() string {
	return "discord"
}

func (d *DiscordClient) Destination() string {
	return d.ChannelID
}
