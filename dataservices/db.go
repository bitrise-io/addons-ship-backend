package dataservices

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	dbDialect string = "postgres"
)

var (
	sqlDB *gorm.DB
)

// ConnectionParams ...
type ConnectionParams struct {
	host     string
	user     string
	dbName   string
	password string
	sslMode  string
}

// GetDB ...
func GetDB() *gorm.DB {
	return sqlDB
}

// Close ...
func Close() {
	closeDB(sqlDB)
	setDB(nil)
}

// InitializeConnection ...
func InitializeConnection(defaultParams ConnectionParams) error {
	connString, err := connectionString(defaultParams)
	if err != nil {
		return errors.WithStack(err)
	}
	db, err := gorm.Open(dbDialect, connString)
	if err != nil {
		return errors.Wrap(err, "Failed to open database")
	}
	setDB(db)
	if err = sqlDB.DB().Ping(); err != nil {
		closeDB(sqlDB)
		return errors.Wrap(err, "Failed to ping database")
	}
	isLogModeEnabled, err := strconv.ParseBool(os.Getenv("GORM_LOG_MODE_ENABLED"))
	if err == nil && isLogModeEnabled {
		sqlDB.LogMode(true)
	}
	sqlDB.LogMode(true)
	return nil
}

func setDB(sdb *gorm.DB) {
	sqlDB = sdb
}

func closeDB(dbToClose *gorm.DB) {
	if dbToClose != nil {
		if err := dbToClose.Close(); err != nil {
			log.Printf(" [!] Exception: Failed to close DB: %+v", err)
		}
	}
}

func (cp ConnectionParams) validate() error {
	if cp.host == "" {
		return errors.New("No database host specified")
	}
	if cp.dbName == "" {
		return errors.New("No database name specified")
	}
	if cp.user == "" {
		return errors.New("No database user specified")
	}
	if cp.password == "" {
		return errors.New("No database password specified")
	}
	return nil
}

func connectionString(defaultParams ConnectionParams) (string, error) {
	connParams := defaultParams
	if connParams.host == "" {
		connParams.host = os.Getenv("DB_HOST")
	}
	if connParams.dbName == "" {
		connParams.dbName = os.Getenv("DB_NAME")
	}
	if connParams.user == "" {
		connParams.user = os.Getenv("DB_USER")
	}
	if connParams.password == "" {
		connParams.password = os.Getenv("DB_PWD")
	}
	if connParams.sslMode == "" {
		connParams.sslMode = os.Getenv("DB_SSL_MODE")
	}
	if err := connParams.validate(); err != nil {
		return "", err
	}
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		connParams.host, connParams.dbName, connParams.user, connParams.password)
	// optionals
	if connParams.sslMode != "" {
		connString += " sslmode=" + connParams.sslMode
	}
	return connString, nil
}
