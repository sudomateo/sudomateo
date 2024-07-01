package main

import (
	"errors"
	"testing"
)

// Compile-time assertion.
var _ Notifier = &MockNotifier{}

type MockNotifier struct {
	notifyCalled bool
	hasErrored   bool
	shouldError  bool
}

func (m *MockNotifier) Notify(message string) error {
	m.notifyCalled = true

	if m.shouldError {
		m.hasErrored = true
		return errors.New("mock error")
	}

	return nil
}

func (m *MockNotifier) Type() string {
	return "mock_notifier"
}

func (m *MockNotifier) Destination() string {
	return "mock_destination"
}

func TestSendNotification(t *testing.T) {
	tt := []struct {
		name string
		m    *MockNotifier
	}{
		{
			name: "happy path",
			m:    &MockNotifier{},
		},
		{
			name: "error path",
			m: &MockNotifier{
				shouldError: true,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			SendNotification("foo", []Notifier{tc.m})

			if !tc.m.notifyCalled {
				t.Fatal("expected notify to be called but wasn't")
			}

			if tc.m.shouldError && !tc.m.hasErrored {
				t.Fatal("wanted error but didn't get one")
			}

			if !tc.m.shouldError && tc.m.hasErrored {
				t.Fatal("did not want error but got one")
			}
		})
	}
}
