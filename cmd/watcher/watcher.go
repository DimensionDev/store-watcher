package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/robfig/cron"
	"github.com/thoas/go-funk"
)

type Callback interface {
	OnUpdate(*ChangedMessage) error
	OnBackupLatestState(map[string]*LatestState)
	OnRestoreLatestState() map[string]*LatestState
}

type Watcher struct {
	latestUpdate map[string]*LatestState
	Callback
}

func (w *Watcher) Start(latestInterval string, targets []*Target) (table *cron.Cron, err error) {
	table = cron.New()
	log.Print("Restore latest state")
	w.latestUpdate = w.OnRestoreLatestState()
	_ = table.AddFunc(latestInterval, func() {
		log.Print("Backup latest state")
		w.OnBackupLatestState(w.latestUpdate)
	})
	for _, item := range targets {
		_ = table.AddFunc(item.Interval, w.makeWatcher(item))
	}
	return
}

func (w *Watcher) makeWatcher(target *Target) func() {
	log.Printf("Add %s to Watcher", target)
	re := regexp.MustCompile(target.Pattern)
	if funk.IndexOfString(re.SubexpNames(), "version") == -1 {
		log.Fatalf("%s Not found `version` regex group in `pattern`", target)
	}
	return func() {
		state := w.getState(target.ID())
		version := w.fetchVersion(target, re)
		if version == "" {
			log.Printf("%s Version info not found", target)
			return
		} else if state.Version == version {
			log.Printf("%s No change", target)
			return
		}
		message := &ChangedMessage{
			Name:            target.Name,
			Platform:        target.Platform,
			Link:            target.Link,
			PreviousVersion: state.Version,
			PreviousDate:    state.Updated,
			CurrentVersion:  version,
			CurrentDate:     time.Now(),
		}
		if message.PreviousDate != nil {
			delta := message.CurrentDate.Sub(*message.PreviousDate).String()
			message.Delta = &delta
		}
		log.Print(message, " Started")
		if err := w.OnUpdate(message); err != nil {
			log.Print(message, err)
		} else {
			state.Updated = &message.CurrentDate
			state.Version = message.CurrentVersion
			log.Print(message, " End")
		}
	}
}

func (w *Watcher) fetchVersion(target *Target, re *regexp.Regexp) (version string) {
	request, _ := http.NewRequest(http.MethodGet, target.Link, nil)
	if target.UserAgent != "" {
		request.Header.Set("User-Agent", target.UserAgent)
	}
	if response, err := http.DefaultClient.Do(request); err != nil {
		return
	} else if response.StatusCode != 200 {
		return
	} else if data, err := ioutil.ReadAll(response.Body); err != nil {
		return
	} else {
		index := funk.IndexOfString(re.SubexpNames(), "version")
		matched := re.FindStringSubmatch(string(data))
		version = matched[index]
	}
	return
}

func (w *Watcher) getState(name string) *LatestState {
	if _, ok := w.latestUpdate[name]; !ok {
		w.latestUpdate[name] = &LatestState{
			Version: "N/A",
			Updated: nil,
		}
	}
	return w.latestUpdate[name]
}
