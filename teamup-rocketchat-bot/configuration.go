package main

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// Configuration holds different options
// required to run the bot
type Configuration struct {
	URL          string `yaml:"URL"`
	Username     string `yaml:"USERNAME"`
	Password     string `yaml:"PASSWORD"`
	UseSSL       bool   `yaml:"USE_SSL"`
	Room         string `yaml:"ROOM"`
	MeetingsCode string `yaml:"MEETINGS_TEAMUP"`
	TeamupToken  string `yaml:"TOKEN_TEAMUP"`
	RepeatIn     int    `yaml:"REPEAT_IN"`
	LogPath      string `yaml:"LOG_PATH"`
	LogFileName  string `yaml:"LOG_FILE_NAME"`
}

// Prints beatiful
func (config Configuration) String() string {
	return fmt.Sprintf(
		"URL:%s\nUsername:%v\nPassword:%s\nUseSSL:%v\nRoom:%s\nMeetingsCode:%s\nTeamupToken:%s\nRepeatIn:%d\nLogPath:%s\nLogFileName:%s\n",
		config.URL,
		config.Username,
		config.Password,
		config.UseSSL,
		config.Room,
		"hidden-for-security-purpose",
		"hidden-for-security-purpose",
		config.RepeatIn,
		config.LogPath,
		config.LogFileName,
	)

}

// IsUrl does basic url checking
// Scheme must not be empty
func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != "" && (u.Scheme == "https" || u.Scheme == "http")
}

// checkValidity checks the validity of the loaded configuration
// and returns error if any field value is not valid
func (config *Configuration) checkValidity(isCustomPath bool) (*Configuration, error) {
	if !isUrl(config.URL) {
		return nil, fmt.Errorf("invalid URL: %s", config.URL)
	}

	if len(config.Username) == 0 {
		return nil, fmt.Errorf("empty USERNAME field")
	}

	if len(config.Password) == 0 {
		return nil, fmt.Errorf("empty PASSWORD field")
	}

	if len(config.Room) == 0 {
		return nil, fmt.Errorf("empty ROOM field")
	}

	if len(config.MeetingsCode) == 0 {
		return nil, fmt.Errorf("empty MEETINGS_TEAMUP field")
	}

	if len(config.TeamupToken) == 0 {
		return nil, fmt.Errorf("empty TOKEN_TEAMUP field")
	}

	if len(config.Room) == 0 {
		return nil, fmt.Errorf("empty ROOM field")
	}
	// if only custom path is not provided
	if !isCustomPath {
		if len(config.LogPath) == 0 {
			// set to default
			config.LogPath = defaultLogFilePath
		} else {
			config.LogPath = path.Clean(config.LogPath) // Clean the path user provided location only
		}
		// Check if the log output folder is accessible
		_, err := os.Stat(config.LogPath)
		if err != nil {
			return nil, fmt.Errorf("the log output folder %s is inaccessible because %s", config.LogPath, err.Error())
		}
	}

	if len(config.LogFileName) == 0 {
		// set to default
		config.LogFileName = defaultLogFileName
	}

	if config.RepeatIn == 0 {
		// set to default
		config.RepeatIn = defaultRepeatIn
	}
	return config, nil
}

// readConfig reads config.yml files
// and returns config or
// corresponding error if any
func readConfig(filePath string, isCustomPath bool) (*Configuration, error) {
	config := Configuration{}

	yamlData, err := os.ReadFile(filePath)

	// Check for errors
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlData, &config)

	if err != nil {
		return nil, err
	}

	validconfig, err := config.checkValidity(isCustomPath)

	if err != nil {
		return &config, err
	}

	return validconfig, nil
}
