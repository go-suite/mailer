package model

type User struct {
	Username       string          `yaml:"username" toml:"username"`
	Password       string          `yaml:"password" toml:"password"`
	Authentication *authentication `yaml:"authentication,omitempty" toml:"authentication,omitempty"`
}
