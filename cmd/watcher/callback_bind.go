package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

type CallbackBind struct {
	*Configure
}

func (c *CallbackBind) OnUpdate(message *ChangedMessage) error {
	log.Print(message, " Started")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.HookProgram)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		//noinspection GoUnhandledErrorResult
		defer stdin.Close()
		encoded, err := json.Marshal(message)
		if err != nil {
			return
		}
		_, _ = stdin.Write(encoded)
	}()
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Print(message, err, "\n", string(output))
		return err
	}
	log.Print(message, " End")
	return nil
}

func (c *CallbackBind) OnBackupLatestState(state map[string]*LatestState) {
	log.Print("Backup latest state")
	data, _ := json.MarshalIndent(state, "", "  ")
	err := ioutil.WriteFile(c.LatestStatePath, data, 0755)
	if err != nil {
		log.Print(err)
	}
}

func (c *CallbackBind) OnRestoreLatestState() (state map[string]*LatestState) {
	log.Print("Restore latest state")
	state = make(map[string]*LatestState)
	data, err := ioutil.ReadFile(c.LatestStatePath)
	if err != nil {
		log.Print(err)
		return
	}
	_ = json.Unmarshal(data, &state)
	return
}
