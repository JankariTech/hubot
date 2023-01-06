### teamup-rocket.chat-bot
- The bot reminds of any upcoming events in teamup calendar(JankariTech's meetings calendar) to the rocket.chat channel
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

