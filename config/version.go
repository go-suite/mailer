package config

var Version string

func init() {
	if len(Version) == 0 {
		Version = "development"
	}
}
