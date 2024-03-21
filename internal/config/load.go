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

	return config, nil
}

// expandEnv iterates through all supported fields and looks for the ENV: prefix
// when located, the environment is checked for the value and if not set, an error is returned
// if set, the value is replaced in the config
func (c *Config) expandEnv() error {
	var errs []error

	c.Root.Username = checkSingleEnvValue(c.Root.Username, &errs)
	c.Root.Password = checkSingleEnvValue(c.Root.Password, &errs)
	c.Defaults.NetworkRestriction = checkSingleEnvValue(c.Defaults.NetworkRestriction, &errs)

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
