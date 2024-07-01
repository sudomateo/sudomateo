package main

import "fmt"

type EmailClient struct {
	Address string
}

func (e *EmailClient) Notify(message string) error {
	if err := simulateError(); err != nil {
		return fmt.Errorf("failed sending email: %w", err)
	}

	return nil
}

func (e *EmailClient) Type() string {
	return "email"
}

func (e *EmailClient) Destination() string {
	return e.Address
}
