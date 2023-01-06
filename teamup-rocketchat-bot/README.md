### teamup rocket.chat bot

- Author: Roshan Lamichhane

### Currently supports

- Reminding upcoming meeting events

### Run locally

*Note*: Make sure that the bot has already been added to channel or room

- Clone the repo and `cd` inside the folder
- Copy `config.yml.example` to `config.yml` and edit the file as required
- Then run the following

```bash
go run main.go
```

**Note**: It will create extra files which are required for smooth functioning of the bot

### Build for cloud or other targets

- Go supports cross-compiling for multiple platforms

```bash
GOARCH=arm GOOS=linux go build -o build/go-jankari-bot 

```
 - The executable will be inside `build` folder

Here,
- `GOARCH` = `targeted architecture` <br>
  E.g. amd64, arm
- `GOOS` = `target operating-system` <br>
  E.g Linux

*Note*: See below for more available options

### Supported combinations

The following table shows the possible combinations of GOOS and GOARCH you can use:
| GOOS | GOARCH | Comment |
| --- | --- | --- |
| android | arm | |
| darwin | 386 | |
| darwin | amd64 | |
| darwin | arm | |
| darwin | arm64 | |
| dragonfly | amd64 | |
| freebsd | 386 | |
| freebsd | amd64 | |
| freebsd | arm | |
| linux | 386 | |
| linux | amd64 | This is the one you need for linux servers and desktops|
| linux | arm | |
| linux | arm64 | |
| linux | ppc64 | |
| linux | ppc64le | |
| linux | mips | |
| linux | mipsle | |
| linux | mips64 | |
| linux | mips64le | |
| netbsd | 386 | |
| netbsd | amd64 | |
| netbsd | arm | |
| openbsd | 386 | |
| openbsd | amd64 | |
| openbsd | arm | |
| plan9 | 386 | |
| plan9 | amd64 | |
| solaris | amd64 | |
| windows | 386 | |
| windows | amd64 | |

### autostart with systemd
1. create a file called `/etc/systemd/system/teamup-rocketchat-bot.service` with this content (remember to adjust the path)

```
[Unit]
Description=teamup reminder for rocket chat systemd service.

[Service]
Type=simple
ExecStart=/bin/bash -c 'cd <path-of-rocketchat-bot>; ./teamup-rocketchat-bot'

[Install]
WantedBy=multi-user.target
```

2. adjust permissions
   `sudo chmod 644 /etc/systemd/system/teamup-rocketchat-bot.service`
3. enable the systemd script
   `sudo systemctl enable teamup-rocketchat-bot.service`
4. start the script
   `systemctl start teamup-rocketchat-bot.service`
6. check if it runs
   `systemctl status teamup-rocketchat-bot.service`
