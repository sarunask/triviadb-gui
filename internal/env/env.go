package env

import (
	"github.com/spf13/pflag"
)

type (
	Config struct {
		TemplatesDir string
	}
)

// Settings holds all settings we have in our app
var Settings *Config

func init() {
	templateDir := pflag.String("template-dir", "templates", "encryption type to be used in S3")
	pflag.Parse()
	Settings = &Config{
		TemplatesDir: *templateDir,
	}
}
