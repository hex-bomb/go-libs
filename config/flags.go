package config

import (
	"github.com/spf13/pflag"
)

var (
	ConfigFileFlag string
)

func init() {
	pflag.StringVarP(&ConfigFileFlag, "config", "c", "", "path to config")
	pflag.Parse()
}
