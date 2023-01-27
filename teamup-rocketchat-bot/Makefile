## Name of the service file
UNIT_FILE_NAME=teamup-rocketchat-bot.service

## Location to keep the unit file
UNIT_FILE_LOCATION=/etc/systemd/system/$(UNIT_FILE_NAME)

## Necessary content for the unit file
LOCAL_LOCATION=./systemd/$(UNIT_FILE_NAME)

SHELL := /bin/bash


## Help information
PHONY: help
help:
	@echo -e "\nSupports following commands"
	@echo -e "=======================================================================================\n"
	@echo -e "build\t\t - Builds the executable and put in the build directory."
	@echo -e "install\t\t - Builds the executable then put into /usr/local/bin and also add a systemd unit file. Requires sudo previleges."
	@echo -e "uninstall\t - Removes the executable from /usr/local/bin and also the systemd unit file. Requires sudo previleges."
	@echo -e "help\t\t - Display this screen."
	@echo -e "\n\nNote: \tMake sure to edit the $(LOCAL_LOCATION) with proper configuration before installing"
PHONY: build
build:
	rm -rf ./build
	@echo "BUILD process started."
	go build -o build/ ./


PHONY: install
install: build
	@echo -e "\n>> COPYING build/* to /usr/local/bin <<"
	sudo cp build/* /usr/local/bin
	@echo -e "\n>> ADDING systemd unit file <<"
	sudo cp $(LOCAL_LOCATION) $(UNIT_FILE_LOCATION)
	sudo chmod 644 $(UNIT_FILE_LOCATION)
	@echo -e "\n>> ENABLING systemd unit file <<"
	sudo systemctl enable $(UNIT_FILE_NAME)
	@echo -e "\n>> TYPE sudo systemctl start $(UNIT_FILE_NAME) to start the bot now. <<"


PHONY: uninstall
uninstall: 
	@echo -e "\n>> REMOVING teamup-rocketchat-bot from /usr/local/bin <<"
	sudo rm /usr/local/bin/teamup-rocketchat-bot
	@echo -e "\n>> STOPPPING the systemd unit service <<"
	sudo systemctl stop $(UNIT_FILE_NAME)
	@echo -e "\n>> DISABLING the systemd unit service <<"
	sudo systemctl disable $(UNIT_FILE_NAME)
	@echo -e "\n>> REMOVING the systemd unit file <<"
	sudo rm $(UNIT_FILE_LOCATION)

