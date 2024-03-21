package config

// Validate does checks purely against the config
// It does not check that all fields are necessarily valid for the databae engine
// but instead checks things like "does a user exist on a database that does not exist in the
// user section of the configuration
func (c *Config) Validate() error {
	// @TODO

	return nil
}
