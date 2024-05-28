package iface

// Manager is the interface
type Manager interface {
	// CreateDatabase ensures that a database exists in the database server
	CreateDatabase(databaseName string) error

	// CreateUser ensures the user exists in the database server
	// Some servers may not utilize the password and/or networkRestriction in this step, but it's here for most flexibility
	CreateUser(username, password, networkRestriction string, globalPerms []string) error

	// AssignWriteUserToDatabase adds a write user to a database
	//
	// This function assumes that CreateUser has been called prior to this function, and that the user already exists
	// if the database engine requires the user to exist
	//
	// Some servers may not utilize the password and/or networkRestriction in this step, but it's here for most flexibility
	AssignWriteUserToDatabase(databaseName, username, password, networkRestriction string) error

	// AssignReadUserToDatabase adds a read-only user to a database
	//
	// This function assumes that CreateUser has been called prior to this function, and that the user already exists
	// if the database engine requires the user to exist
	//
	// Some servers may not utilize the password and/or networkRestriction in this step, but it's here for most flexibility
	AssignReadUserToDatabase(databaseName, username, password, networkRestriction string) error
}
