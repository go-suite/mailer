package config

import (
	"github.com/gennesseaux/mailer/model"
	"github.com/go-mods/convert"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"os"
)

var C *Config

type Config struct {
	Users  []*model.User
	Debug  bool
	Secure bool
}

func init() {
	// Create an instance of the configuration struct
	C = &Config{}
	C.Debug, _ = convert.ToBool(os.Getenv("MAILER_DEBUG"))
	C.Secure = false

	// Feed the configuration data from a TOML file
	if _, err := os.Stat("data/users.toml"); err == nil {
		// Create an instance of the configuration struct
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
		// Create an instance of the configuration struct
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
