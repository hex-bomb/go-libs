package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hex-bomb/go-libs/validator"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var (
	// appName is an application build name
	appName = "unknown"
	// appVersion is an application build version.
	appVersion = "development"
)

// AppName returns an application build name
func AppName() string {
	return appName
}

// AppVersion returns an application build version.
func AppVersion() string {
	return appVersion
}

var hooks []viper.DecoderConfigOption

// RegisterHook adds hooks to configs unmarshalling process.
func RegisterHook(hook mapstructure.DecodeHookFuncType) {
	hooks = append(hooks, viper.DecodeHook(hook))
}

// GetSettings loads configuration from environmental variables to settings.
func GetSettings[AppSettings any](defaults map[string]any) (*AppSettings, error) {
	var settings AppSettings
	var err error

	if err = configViper(defaults); err != nil {
		return nil, err
	}

	if err = readConfigs(&settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// configViper sets default viper configuration.
func configViper(defaults map[string]any, paths ...string) error {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath("./")
	viper.AddConfigPath("./configs")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}

	viper.SetConfigName("config")
	if ConfigFileFlag != "" {
		if _, err := os.Stat(ConfigFileFlag); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("specified config path '%s' is not exists", ConfigFileFlag)
		}

		viper.SetConfigFile(ConfigFileFlag)
	}

	viper.AutomaticEnv()
	return nil
}

// readConfigs reads env variables, validates and create configs.
func readConfigs(settings any) error {
	if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return err
	}

	if err := viper.Unmarshal(settings, hooks...); err != nil {
		return fmt.Errorf("could not unmarshal configuration: %w", err)
	}

	if err := validator.Get().Struct(settings); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
