### teamup-rocket.chat-bot
- The bot reminds of any upcoming events in teamup calendar(JankariTech's meetings calendar) to the rocket.chat channel
- Author: Roshan Lamichhane

### Currently supports

- Reminding upcoming meeting events
### Configuration
- The configuration file should have the `yml` extension and the following properties.
```bash
# Url of rocket chat server
URL: https://rocket.chat

# bot username
USERNAME: bot

# bot password
PASSWORD: pass

# Use ssl(https)?
USE_SSL: TRUE

# channel/room (Currently supports one)
# Make sure the bot has already been added to the channel
ROOM: try

# meeting calendar read-only link (calendar-code)
MEETINGS_TEAMUP: someMeetingCodexxxxx

# teamup token required for api calls
TOKEN_TEAMUP: someLongCodexxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# run the bot in every X interval (mins)
# Keep it between 5 to 10
# default is 5
REPEAT_IN: 5

# path to keep logfile
# #default LOG_PATH is /tmp folder
# It will be overriden if logpath flag is provided to bot while running the bot
LOG_PATH: /tmp

# name for the log file
# default log file name is teamup-rocket-chat.log
LOG_FILE_NAME: teamup-rocket-chat.log
```
- The application supports parameters for custom logpath and custom configuraion file path
- `--config Point to the configuration (config) file. It overrides the default configuration file located at app directory`
- `--logpath Set the custom logpath. It overrides the log path specified in configuration (config.yml) file`
- `--help Show help and usage screen`

#### templates
Messages can be customized by calendar, for that Go-templates are used.
1. get the ID of your calendar:
   1. go to the settings page of teamup
   2. go to the calendars list
   3. click on the calendar you want to create a template for
   4. copy the number at the end of the URL e.g. `123456` of `https://teamup.com/c/xdfesf/jankaritech/settings/calendars/edit/123456`, this is your calendar id
2. copy `example.tmpl` to `<calendar-id>.tmpl`
3. [adjust the template](https://pkg.go.dev/text/template). As functions all [sprig](https://masterminds.github.io/sprig/) functions can be used plus `NotesInMarkdown`, which will return the notes converted from HTML to MD.
4. repeat for every calendar for which you want to have a custom message for
5. if there is no template found a build in, default template will be used.

### Run locally

*Note*: Make sure that the bot has already been added to channel or room

- Clone the repo and `cd` inside the folder
- Copy `config.yml.example` to `config.yml` and edit the file as required
- Then run the following

```bash
go run .
```

- For more application usage, run 
```bash
go run . --help
```

**Note**: It will create extra files which are required for smooth functioning of the bot. By default, log file will be created at `/tmp` folder with name `teamup-rocketchat.log`

### Build for cloud or other targets

- Go supports cross-compiling for multiple platforms

```bash
GOARCH=amd64 GOOS=linux go build -o build/ .
```
 - The executable will be inside `build` folder

Here,
- `GOARCH` = `targeted architecture` <br>
  E.g. amd64, arm
- `GOOS` = `targeted operating-system` <br>
  E.g. linux,windows

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
