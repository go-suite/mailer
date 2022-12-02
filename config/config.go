package config

import (
	"github.com/go-mods/convert"
	"github.com/go-suite/mailer/model"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"os"
)

var C *Config

type Config struct {
	// Users Array of users allowed to send mail.
	// Each user must authenticate with a login password.
	Users []*model.User `yaml:"users" toml:"users"`

	// Secure is set to true if a least one user is defined.
	// In secure mode, only authenticated users can send mails.
	Secure bool

	// Debug if true, the router will be more verbose.
	Debug bool

	// LogFile holds the configuration used for the log file
	LogFile logFile `yaml:"logfile"`
}

type logFile struct {
	// Filename is the file to write logs to. Backup log files will be retained
	// in the same directory. It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string `yaml:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `yaml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `yaml:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `yaml:"maxbackups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `yaml:"localtime"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `yaml:"compress"`
	// contains filtered or unexported fields
}

func init() {
	// Create an instance of Config
	C = &Config{}

	// Read the MAILER_DEBUG environment variable
	C.Debug, _ = convert.ToBool(os.Getenv("MAILER_DEBUG"))

	// The server is not secured at startup
	C.Secure = false

	// Initialise users list
	initUsers()

	// Initialise logfile parameters
	initLogFile()
}

func initUsers() {
	// Feed the configuration data from a TOML file
	if _, err := os.Stat("data/users.toml"); err == nil {
		// Create an instance of Config
		tomlConfig := &Config{}
		// Create a feeder that provides the configuration data from a TOML file
		tomlFeeder := feeder.Toml{Path: "data/users.toml"}
		// Create a Config instance and feed `tomlConfig` using `tomlFeeder`
		c := config.New()
		c.AddFeeder(tomlFeeder)
		c.AddStruct(tomlConfig)
		err := c.Feed()
		if err == nil {
			for _, u := range tomlConfig.Users {
				C.addUser(u)
			}
			C.Secure = true
		}
	}

	// Feed the configuration data from a YAML file
	if _, err := os.Stat("data/users.yaml"); err == nil {
		// Create an instance of Config
		yamlConfig := &Config{}
		// Create a feeder that provides the configuration data from a YAML file
		yamlFeeder := feeder.Yaml{Path: "data/users.yaml"}
		// Create a Config instance and feed `yamlConfig` using `yamlFeeder`
		c := config.New()
		c.AddFeeder(yamlFeeder)
		c.AddStruct(yamlConfig)
		err := c.Feed()
		if err == nil {
			for _, u := range yamlConfig.Users {
				C.addUser(u)
			}
			C.Secure = true
		}
	}
}

func initLogFile() {
	// Create an instance of Config LogFile
	yamlConfig := &Config{
		LogFile: logFile{
			Filename: "mailer.log",
			MaxSize:  10,
		},
	}

	if _, err := os.Stat("data/log.yaml"); err == nil {
		// Create a feeder that provides the configuration data from a YAML file
		yamlFeeder := feeder.Yaml{Path: "data/log.yaml"}
		// Create a LogFile instance and feed `logFile` using `yamlFeeder`
		c := config.New()
		c.AddFeeder(yamlFeeder)
		c.AddStruct(yamlConfig)
		err := c.Feed()
		if err != nil {
			panic("Error while loading log.yaml")
		} else {
			C.LogFile = yamlConfig.LogFile
		}
	}
}

// addUser adds a user if it does not exist yet
func (c *Config) addUser(user *model.User) {
	for _, u := range c.Users {
		if u.Username == user.Username {
			return
		}
	}
	c.Users = append(c.Users, user)
}

// GetUser gets a user from name
func (c *Config) GetUser(username string) *model.User {
	for _, u := range c.Users {
		if u.Username == username {
			return u
		}
	}
	return nil
}
