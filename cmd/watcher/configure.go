package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configure struct {
	Interval            string    `json:"interval,omitempty"`
	UserAgent           string    `json:"user_agent,omitempty"`
	LatestStateInterval string    `json:"latest_state_interval,omitempty"`
	LatestStatePath     string    `json:"latest_state_path,omitempty"`
	HookProgram         string    `json:"hook_program,omitempty"`
	Watches             []*Target `json:"watches"`
}

func loadConfigure(name string) (*Configure, error) {
	configure := &Configure{
		Interval:            "0 0 * * * *",  // at minute 0
		LatestStateInterval: "0 30 * * * *", // at minute 30
		UserAgent:           "store-watcher/1.0 powered-by dimension.im",
		LatestStatePath:     "./latest-state.json",
		HookProgram:         "./on-update.py",
	}
	if fp, err := os.Open(name); err != nil {
		return nil, err
	} else if data, err := ioutil.ReadAll(fp); err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, configure); err != nil {
		return nil, err
	}
	for _, watch := range configure.Watches {
		if watch.Interval == "" {
			watch.Interval = configure.Interval
		}
		if watch.UserAgent == "" {
			watch.UserAgent = configure.UserAgent
		}
	}
	return configure, nil
}
