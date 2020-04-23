package main

import (
	"fmt"
	"time"
)

type Target struct {
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	Link      string `json:"target"`
	Pattern   string `json:"pattern"`
	Interval  string `json:"interval,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

func (t *Target) ID() string {
	return fmt.Sprintf("%s|%s", t.Name, t.Platform)
}

func (t *Target) String() string {
	return fmt.Sprintf("[%s - %s]", t.Name, t.Platform)
}

type LatestState struct {
	Version string     `json:"version"`
	Updated *time.Time `json:"updated"`
}

type ChangedMessage struct {
	Name            string     `json:"name"`
	Platform        string     `json:"platform"`
	Link            string     `json:"link"`
	PreviousVersion string     `json:"previous_version"`
	PreviousDate    *time.Time `json:"previous_date"`
	CurrentVersion  string     `json:"current_version"`
	CurrentDate     time.Time  `json:"current_date"`
}

func (m *ChangedMessage) String() string {
	return fmt.Sprintf("[%s - %s] %s -> %s", m.Name, m.Platform, m.PreviousVersion, m.CurrentVersion)
}
