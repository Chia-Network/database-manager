package config

import (
	"fmt"
)

// Config is the parent config object
type Config struct {
	Connection Connection `yaml:"connection"`
	Defaults   Defaults   `yaml:"defaults"`
	Users      []User     `yaml:"users"`
	Databases  []Database `yaml:"databases"`

	userLookup map[string]*User
}

// Connection is the user/connection details that has permissions to create databases and users
// It doesn't have to actually be the root user, but it must have appropriate permissions
type Connection struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

// Defaults for all users and/or databases can be defined here
type Defaults struct {
	NetworkRestrictions []string `yaml:"network_restrictions"`
}

// User is a single database user
// Permissions are defined per DB
type User struct {
	Username          string   `yaml:"username"`
	Password          string   `yaml:"password"`
	GlobalPermissions []string `yaml:"global-permissions"`
}

// Database is a database + its user permissions
// Users and ReadonlyUsers must be defined in the users section of the config of this will fail
type Database struct {
	Name          string   `yaml:"name"`
	Users         []string `yaml:"users"`
	ReadonlyUsers []string `yaml:"readonly_users"`
}

// GetUserByUsername returns a user by the username
func (c *Config) GetUserByUsername(username string) (*User, error) {
	user, exists := c.userLookup[username]
	if !exists {
		return nil, fmt.Errorf("can't find user in config")
	}

	return user, nil
}
