package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads config from the specified filename
func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s", configPath)
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}

	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config yaml: %w", err)
	}

	err = config.expandEnv()
	if err != nil {
		return nil, err
	}

	config.userLookup = map[string]*User{}
	for _, user := range config.Users {
		config.userLookup[user.Username] = &user
	}

	return config, nil
}

// expandEnv iterates through all supported fields and looks for the ENV: prefix
// when located, the environment is checked for the value and if not set, an error is returned
// if set, the value is replaced in the config
func (c *Config) expandEnv() error {
	var errs []error

	c.Connection.Username = checkSingleEnvValue(c.Connection.Username, &errs)
	c.Connection.Password = checkSingleEnvValue(c.Connection.Password, &errs)
	c.Connection.Host = checkSingleEnvValue(c.Connection.Host, &errs)
	// @TODO deal with string vs uint
	//c.Connection.Port = checkSingleEnvValue(c.Connection.Port, &errs)

	for i, restriction := range c.Defaults.NetworkRestrictions {
		c.Defaults.NetworkRestrictions[i] = checkSingleEnvValue(restriction, &errs)
	}
	// If empty, default to localhost (safer than %)
	if len(c.Defaults.NetworkRestrictions) == 0 {
		c.Defaults.NetworkRestrictions = append(c.Defaults.NetworkRestrictions, "localhost")
	}

	for i, user := range c.Users {
		c.Users[i].Username = checkSingleEnvValue(user.Username, &errs)
		c.Users[i].Password = checkSingleEnvValue(user.Password, &errs)
	}

	for i, database := range c.Databases {
		c.Databases[i].Name = checkSingleEnvValue(database.Name, &errs)

		for d, user := range database.Users {
			c.Databases[i].Users[d] = checkSingleEnvValue(user, &errs)
		}
		for d, user := range database.ReadonlyUsers {
			c.Databases[i].ReadonlyUsers[d] = checkSingleEnvValue(user, &errs)
		}
	}

	if len(errs) > 0 {
		errStr := fmt.Sprintf("%d errors encountered processing config\n", len(errs))
		for _, err := range errs {
			errStr = fmt.Sprintf("%s%s\n", errStr, err.Error())
		}
		return fmt.Errorf("%s", errStr)
	}

	return nil
}

func checkSingleEnvValue(value string, errs *[]error) string {
	if !strings.HasPrefix(value, "ENV:") {
		return value
	}

	envKey := value[4:]
	envValue, isSet := os.LookupEnv(envKey)
	if !isSet {
		*errs = append(*errs, fmt.Errorf("environment variable %s is not set", envKey))
		return ""
	}

	return envValue
}
