package mocks

import (
	"github.com/brienze1/notes-api/internal/integration/persistence/repositories/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var once sync.Once
var db *gorm.DB

// Db class is used to configure DB and create a connection pool using gorm
func Db() *gorm.DB {
	if db == nil {
		once.Do(
			func() {
				db = open()
			},
		)
	}

	return db
}

func open() *gorm.DB {
	dbConn, err := gorm.Open(
		sqlite.Open(":memory:"), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)
	if err != nil {
		panic("failed to connect to database. err: " + err.Error())
	}

	err = ClearDB(dbConn)
	for err != nil {
		err = ClearDB(dbConn)
	}

	return dbConn
}

func ClearDB(db *gorm.DB) error {
	db.Exec("ATTACH ':memory:' AS public")

	if err := db.Migrator().DropTable(
		&models.User{},
		&models.Note{},
	); err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Note{},
	); err != nil {
		return err
	}

	return nil
}
