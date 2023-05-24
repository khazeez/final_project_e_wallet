package app

import (
	"errors"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
)

type UserUsecase interface {
	InsertUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	FindOne(id int) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Delete(id int) error
}
type userUsecase struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (u *userUsecase) InsertUser(user *domain.User) error {

	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	return u.userRepository.Create(user)
}

func (u *userUsecase) UpdateUser(user *domain.User) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	return u.userRepository.Update(user)
}

func (u *userUsecase) FindOne(id int) (*domain.User, error) {
	return u.userRepository.FindOne(id)
}

func (u *userUsecase) FindAll() ([]domain.User, error) {
	return u.userRepository.FindAll()
}

func (u *userUsecase) Delete(id int) error {
	return u.userRepository.Delete(id)
}
