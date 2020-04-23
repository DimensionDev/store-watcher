package main

var configure *Configure

func init() {
	var err error
	configure, err = loadConfigure("configure.json")
	if err != nil {
		return
	}
}

func main() {
	callback := &CallbackBind{Configure: configure}
	service := &Watcher{Callback: callback}
	cron, err := service.Start(configure.LatestStateInterval, configure.Watches)
	if err != nil {
		panic(err)
	}
	cron.Run()
}
