package iface

// Manager is the interface
type Manager interface {
	// CreateDatabase ensures that a database exists in the database server
	CreateDatabase(databaseName string) error

	// CreateUser ensures the user exists in the database server
	// Some servers may not utilize the password in this step, but it's here for most flexibility
	CreateUser(username, password string) error

	// AssignWriteUserToDatabase adds a write user to a database
	// Some servers may not utilize the password in this step, but it's here for most flexibility
	AssignWriteUserToDatabase(databaseName, username, password string) error

	// AssignReadUserToDatabase adds a read-only user to a database
	// Some servers may not utilize the password in this step, but it's here for most flexibility
	AssignReadUserToDatabase(databaseName, username, password string) error
}
