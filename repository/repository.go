package repository

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository interface {
	CreateRepositoryWithContext(context.Context) Repository
	UserRepository() UserRepository
}

func NewRepository() Repository {
	logger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             800 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=procal password=password dbname=postgres port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		panic("unable to connect to database: " + err.Error())
	}
	return buildRepositoryStruct(db)
}

func (repo *repository) CreateRepositoryWithContext(context context.Context) Repository {
	db := repo.connection.WithContext(context)
	return buildRepositoryStruct(db)
}

func buildRepositoryStruct(db *gorm.DB) Repository {
	return &repository{
		connection:     db,
		userRepository: NewUserRepository(db),
	}
}

type repository struct {
	connection     *gorm.DB
	userRepository UserRepository
}

func (r *repository) UserRepository() UserRepository {
	return r.userRepository
}
