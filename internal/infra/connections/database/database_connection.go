package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brienze1/notes-api/internal/infra/properties"
	"github.com/brienze1/notes-api/internal/integration/adapters"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb(cache adapters.Cache, secret adapters.Secret) *gorm.DB {
	ctx := context.Background()
	secretBytes, err := cache.Get(ctx, properties.CacheDbSecretKey)
	if err != nil || secretBytes == nil {
		secretBytes, err = secret.Get(ctx, properties.GetSecretDatabase())
		if err != nil {
			panic(err)
		}
		err = cache.Set(ctx, properties.CacheDbSecretKey, secretBytes, properties.CacheExpiration)
		if err != nil {
			fmt.Printf("Error setting cache: %s \n", err.Error())
		}
	}

	var secretData map[string]any
	err = json.Unmarshal(secretBytes, &secretData)
	if err != nil {
		fmt.Printf("Error parsing secret: %s \n", err.Error())
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=%s connect_timeout=%s application_name=%s",
		secretData["host"], secretData["username"],
		secretData["password"], secretData["dbname"],
		secretData["port"],
		properties.GetDatabaseSslMode(),
		properties.GetDatabaseTimeout(),
		properties.ApplicationRepositoryName)

	cfg := &gorm.Config{}
	dialect := postgres.Open(dsn)

	db, err := gorm.Open(dialect, cfg)
	if err != nil {
		panic("failed to connect to database. err: " + err.Error())
	}

	database, err := db.DB()
	if err != nil {
		panic("failed to connect to database. err: " + err.Error())
	}

	database.SetMaxOpenConns(properties.GetDatabaseMaxOpenConnections())
	database.SetMaxIdleConns(properties.GetDatabaseMaxIdleConnections())
	return db
}
