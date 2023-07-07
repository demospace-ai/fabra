package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.fabra.io/server/common/application"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/secret"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DB_PASSWORD_KEY = "projects/932264813910/secrets/fabra-db-password/versions/latest"

func InitDatabase() (*gorm.DB, error) {
	if application.IsProd() {
		return initDatabaseProd()
	} else {
		return initDatabaseDev()
	}
}

func initDatabaseDev() (*gorm.DB, error) {
	dbURI := "user=fabra password=fabra database=fabra host=localhost sslmode=require"

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "(database.initDatabaseDev) sql.Open")
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "(database.initDatabaseDev) error getting raw DB handle")
	}

	sqldb.SetMaxOpenConns(10)
	sqldb.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func initDatabaseProd() (*gorm.DB, error) {
	var (
		dbUser = mustGetenv("DB_USER")
		dbHost = mustGetenv("DB_HOST")
		dbPort = mustGetenv("DB_PORT")
		dbName = mustGetenv("DB_NAME")
	)

	dbPwd, err := secret.FetchSecret(context.TODO(), DB_PASSWORD_KEY)
	if err != nil {
		return nil, errors.Wrap(err, "(database.initDatabaseProd) fetching secret)")
	}

	// TODO: use client certificates here and enforce SSL verify full on Cloud SQL
	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s", dbHost, dbUser, *dbPwd, dbPort, dbName)

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, errors.Wrap(err, "(database.initDatabaseProd) sql.Open")
	}

	return db, nil
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Error: %s environment variable not set.\n", k)
	}
	return v
}
