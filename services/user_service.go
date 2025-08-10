package services

import (
	"context"

	"procal/entity"
	"procal/repository"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	FindById(id int) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
	FindByFirebaseUid(firebaseUid string) (entity.User, error)
	Create(user entity.User) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	Delete(id int) error
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) FindById(id int) (entity.User, error) {
	return service.userRepository.FindById(context.Background(), id)
}

func (service *userService) FindByEmail(email string) (entity.User, error) {
	return service.userRepository.FindByEmail(context.Background(), email)
}

func (service *userService) FindByFirebaseUid(firebaseUid string) (entity.User, error) {
	return service.userRepository.FindByFirebaseUid(context.Background(), firebaseUid)
}

func (service *userService) Create(user entity.User) (entity.User, error) {
	return service.userRepository.Create(context.Background(), user)
}

func (service *userService) Update(user entity.User) (entity.User, error) {
	return service.userRepository.Update(context.Background(), user)
}

func (service *userService) Delete(id int) error {
	return service.userRepository.Delete(context.Background(), id)
}
