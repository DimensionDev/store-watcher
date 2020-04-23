# Store Watcher

This tool is used to observe the actual release time of each store

## Build

```bash
go build ./cmd/watcher
```

## File tree

```plain
├── README.md
├── configure.json     # The configure file
├── latest-state.json  # record latest version and updated at
├── on-update.py       # Observed updated hook
└── watcher            # Main program
```

## configure.json [example]

```plain
{
  "interval": "0 * * * *", // at minute 0, fetch all store version info
  "user_agent": "store-watcher/1.0 powered by dimension.im", // request default user-agent
  "latest_state_interval": "30 * * * *", // at minute 30, save latest state to file
  "latest_state_path": "latest-state.json", // latest state save to file path
  "hook_program": "./on-update.py", // observed updated hook program
  "watches": [
    {
      "name": "Maskbook", // product name
      "platform": "Chrome Store", // publish platform
      "target": "https://chrome.google.com/webstore/detail/maskbook/jkoeaghipilijlahjplgbfiocjhldnap",
      "pattern": "version\" content=\"(?P<version>[\\d\\.]+)\"" // version info (regex pattern)
    }
  ]
}
```
