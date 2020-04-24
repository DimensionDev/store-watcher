package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

type CallbackBind struct {
	*Configure
}

func (c *CallbackBind) OnUpdate(message *ChangedMessage) error {
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
		encoded, _ := json.Marshal(message)
		_, _ = stdin.Write(encoded)
	}()
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w\nOutput: %s", err, string(output))
	}
	return nil
}

func (c *CallbackBind) OnBackupLatestState(state map[string]*LatestState) {
	data, _ := json.MarshalIndent(state, "", "  ")
	err := ioutil.WriteFile(c.LatestStatePath, data, 0755)
	if err != nil {
		log.Print(err)
	}
}

func (c *CallbackBind) OnRestoreLatestState() (state map[string]*LatestState) {
	state = make(map[string]*LatestState)
	data, err := ioutil.ReadFile(c.LatestStatePath)
	if err != nil {
		log.Print(err)
		return
	}
	_ = json.Unmarshal(data, &state)
	return
}
