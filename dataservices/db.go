package dataservices

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// Postgres driver
	_ "github.com/lib/pq"
)

const (
	dbDialect string = "postgres"
)

var (
	sqlDB *gorm.DB
)

// ConnectionParams ...
type ConnectionParams struct {
	Host     string
	User     string
	DBName   string
	Password string
	SSLMode  string
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
func InitializeConnection(defaultParams ConnectionParams, withDB bool) error {
	connString, err := connectionString(defaultParams, withDB)
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
	if cp.Host == "" {
		return errors.New("No database host specified")
	}
	if cp.DBName == "" {
		return errors.New("No database name specified")
	}
	if cp.User == "" {
		return errors.New("No database user specified")
	}
	if cp.Password == "" {
		return errors.New("No database password specified")
	}
	return nil
}

func connectionString(defaultParams ConnectionParams, withDB bool) (string, error) {
	connParams := defaultParams
	if connParams.Host == "" {
		connParams.Host = os.Getenv("DB_HOST")
	}
	if connParams.DBName == "" {
		connParams.DBName = os.Getenv("DB_NAME")
	}
	if connParams.User == "" {
		connParams.User = os.Getenv("DB_USER")
	}
	if connParams.Password == "" {
		connParams.Password = os.Getenv("DB_PWD")
	}
	if connParams.SSLMode == "" {
		connParams.SSLMode = os.Getenv("DB_SSL_MODE")
	}
	if err := connParams.validate(); err != nil {
		return "", err
	}
	connString := fmt.Sprintf("host=%s user=%s password=%s",
		connParams.Host, connParams.User, connParams.Password)
	if withDB {
		connString += " dbname=" + connParams.DBName
	}
	// optionals
	if connParams.SSLMode != "" {
		connString += " sslmode=" + connParams.SSLMode
	}
	return connString, nil
}
