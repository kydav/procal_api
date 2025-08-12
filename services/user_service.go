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
	FindById(ctx context.Context, id string) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindByFirebaseUid(ctx context.Context, firebaseUid string) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id string) error
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) FindById(ctx context.Context, id string) (entity.User, error) {
	return service.userRepository.FindById(ctx, id)
}

func (service *userService) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	return service.userRepository.FindByEmail(ctx, email)
}

func (service *userService) FindByFirebaseUid(ctx context.Context, firebaseUid string) (entity.User, error) {
	return service.userRepository.FindByFirebaseUid(ctx, firebaseUid)
}

func (service *userService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	return service.userRepository.Create(ctx, user)
}

func (service *userService) Update(ctx context.Context, user entity.User) (entity.User, error) {
	repoUser, err := service.userRepository.FindById(ctx, user.ID)
	if err != nil {
		return entity.User{}, err
	}
	repoUser.FirstName = user.FirstName
	repoUser.LastName = user.LastName
	repoUser.Email = user.Email
	repoUser.Age = user.Age
	repoUser.Height = user.Height
	repoUser.CurrentWeight = user.CurrentWeight
	repoUser.Gender = user.Gender
	repoUser.MeasurementPreference = user.MeasurementPreference

	return service.userRepository.Update(ctx, repoUser)
}

func (service *userService) Delete(ctx context.Context, id string) error {
	return service.userRepository.Delete(ctx, id)
}
