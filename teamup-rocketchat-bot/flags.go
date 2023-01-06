package main

import (
	"fmt"
)

// banner for the app
const bannerText = `
░░░░▒█░█▀▀▄░█▀▀▄░█░▄░█▀▀▄░█▀▀▄░░▀░░▀▀█▀▀░█▀▀░█▀▄░█░░░
░░░░▒█░█▄▄█░█░▒█░█▀▄░█▄▄█░█▄▄▀░░█▀░░▒█░░░█▀▀░█░░░█▀▀█
░▒█▄▄█░▀░░▀░▀░░▀░▀░▀░▀░░▀░▀░▀▀░▀▀▀░░▒█░░░▀▀▀░▀▀▀░▀░░▀
`
const about = "About:\tA bot that fetches events from teamup calendar and notifies gophers of JankariTech through rocket.chat"

// displayHelpMessage displays help message
func displayHelpMessage(bannerText string) {

	msg := "Usage: teamup-rocketchat-bot --option <value>\n\n"
	msg += "Options:\n"
	msg += "-h, --help \t Show this screen\n"
	msg += "--logpath\t Set the custom logpath. It overrides the log path specified in configuration (config.yml) file\n"
	msg += "--config\t Point to the configuration (config) file. It overrides the default configuration file located at app directory\n"
	msg += "\nExample:\n"
	msg += "teamup-rocketchat-bot --config /home/user/.config/bot/config.yml\n\n\tHere, the bot will read configuration from the pointed file\n\n"
	msg += "teamup-rocketchat-bot --logpath /tmp/\n\n\tHere, the configuration (config.yml) will be read from default path (i.e the app directory) and\n\t the log path provided will override the path specified in configuration file\n\n"
	msg += "teamup-rocketchat-bot --config /home/user/.config/bot/config.yml --logpath /tmp/\n\n\tHere, the configuration (config.yml) will be read from the pointed file and\n\t the log path provided will override the path specified in configuration file\n\n"
	msg += "\nNote:\tMake sure the app has necessary premissions required to read or write at custom locations provided\n"
	msg += "\tAlso, trying to run the bot without providing config.yml at app directory and without providing --config switch will result in failure of application start.\n"
	fmt.Printf("%s\n%s\n%s\n", bannerText, about, msg)
}

func displayHelpMessageWithoutBanner() {
	displayHelpMessage("")
}
