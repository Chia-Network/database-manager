package config

import (
	"fmt"
)

// Validate does checks purely against the config
// It does not check that all fields are necessarily valid for the databae engine
// but instead checks things like "does a user exist on a database that does not exist in the
// user section of the configuration
func (c *Config) Validate() error {
	users := map[string]uint{}
	for _, user := range c.Users {
		if _, set := users[user.Username]; !set {
			users[user.Username] = 0
		}
		users[user.Username]++
		if users[user.Username] > 1 {
			return fmt.Errorf("user %s is in config more than once", user.Username)
		}
	}

	for _, database := range c.Databases {
		for _, expectedUsername := range database.Users {
			if _, userSet := users[expectedUsername]; !userSet {
				return fmt.Errorf("user %s for database %s is not defined in the users section of the config", expectedUsername, database.Name)
			}
		}
		for _, expectedUsername := range database.ReadonlyUsers {
			if _, userSet := users[expectedUsername]; !userSet {
				return fmt.Errorf("user %s for database %s is not defined in the users section of the config", expectedUsername, database.Name)
			}
		}
	}

	return nil
}
