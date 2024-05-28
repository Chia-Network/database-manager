package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/chia-network/database-manager/internal/utils"
)

// Manager implements the manager interface for MySQL databases
type Manager struct {
	client *sql.DB
}

// NewMySQLManager returns a new instance of the mysql manager
func NewMySQLManager(rootUser, rootPassword, host, port string) (*Manager, error) {
	cfg := mysql.Config{
		User:                 rootUser,
		Passwd:               rootPassword,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", host, port),
		AllowNativePasswords: true,
	}
	client, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return &Manager{
		client: client,
	}, nil
}

// CreateDatabase ensures that a database exists in the database server
func (m *Manager) CreateDatabase(databaseName string) error {
	// Can't use ? placeholders in a prepared statement for the DB name
	databaseName, err := sanitizeDatabaseName(databaseName)
	if err != nil {
		return err
	}
	result, err := m.client.Query(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", databaseName))
	if err != nil {
		return err
	}
	utils.LogErrorAndContinue(result.Close(), "error closing result in CreateDatabase")
	return nil
}

// CreateUser ensures the user exists in the database server
func (m *Manager) CreateUser(username, password, networkRestriction string, globalPerms []string) error {
	username, err := sanitizeUsername(username)
	if err != nil {
		return err
	}
	password, err = sanitizePassword(password)
	if err != nil {
		return err
	}
	result, err := m.client.Query(fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%s' IDENTIFIED BY '%s';", username, networkRestriction, password))
	if err != nil {
		return err
	}
	utils.LogErrorAndContinue(result.Close(), "error closing result for create user in CreateUser")
	if len(globalPerms) > 0 {
		result, err := m.client.Query(fmt.Sprintf("GRANT %s on *.* to '%s'@'%s';", strings.Join(globalPerms, ","), username, networkRestriction))
		if err != nil {
			return err
		}
		utils.LogErrorAndContinue(result.Close(), "error closing result for global perms in CreateUser")
	}
	return nil
}

// AssignWriteUserToDatabase adds a write user to a database
//
// This function assumes that CreateUser has been called prior to this function, and that the user already exists
// if the database engine requires the user to exist
func (m *Manager) AssignWriteUserToDatabase(databaseName, username, password, networkRestriction string) error {
	databaseName, err := sanitizeDatabaseName(databaseName)
	if err != nil {
		return err
	}
	username, err = sanitizeUsername(username)
	if err != nil {
		return err
	}
	result, err := m.client.Query(fmt.Sprintf(
		"GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'%s';",
		databaseName,
		username,
		networkRestriction,
	))
	if err != nil {
		return err
	}
	utils.LogErrorAndContinue(result.Close(), "error closing result in AssignWriteUserToDatabase")
	return nil
}

// AssignReadUserToDatabase adds a read-only user to a database
//
// This function assumes that CreateUser has been called prior to this function, and that the user already exists
// if the database engine requires the user to exist
func (m *Manager) AssignReadUserToDatabase(databaseName, username, password, networkRestriction string) error {
	databaseName, err := sanitizeDatabaseName(databaseName)
	if err != nil {
		return err
	}
	username, err = sanitizeUsername(username)
	if err != nil {
		return err
	}
	result, err := m.client.Query(fmt.Sprintf(
		"GRANT SELECT, SHOW VIEW ON `%s`.* TO '%s'@'%s'",
		databaseName,
		username,
		networkRestriction,
	))
	if err != nil {
		return err
	}
	utils.LogErrorAndContinue(result.Close(), "error closing result in AssignReadUserToDatabase")
	return nil
}
