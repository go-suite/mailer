package config

var Version string

func init() {
	// By default, the version is set to development
	if len(Version) == 0 {
		Version = "development"
	}
}
