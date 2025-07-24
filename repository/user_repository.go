package repository

import (
	"context"
	"procal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindById(ctx context.Context, id int) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id int) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		panic(err)
	}
	return &userRepository{connection: db}
}

type userRepository struct {
	connection *gorm.DB
}

// Create implements UserRepository.
func (db *userRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	if result := db.connection.Create(&user); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// Delete implements UserRepository.
func (db *userRepository) Delete(ctx context.Context, id int) error {
	if result := db.connection.Delete(&entity.User{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByEmail implements UserRepository.
func (db *userRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	if result := db.connection.Where("email = ?", email).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// FindById implements UserRepository.
func (db *userRepository) FindById(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	if result := db.connection.First(&user, id); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// Update implements UserRepository.
func (db *userRepository) Update(ctx context.Context, user entity.User) (entity.User, error) {
	if result := db.connection.Save(&user); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
