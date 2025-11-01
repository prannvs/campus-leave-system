package services

import (
	"github.com/prannvs/campus-leave-system/internal/models"
	"github.com/prannvs/campus-leave-system/internal/repositories"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req models.RegisterRequest) (*models.User, error) {
	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, gorm.ErrDuplicatedKey
	}

	user := &models.User{
		Name:   req.Name,
		Email:  req.Email,
		Role:   req.Role,
		Dept:   req.Dept,
		Hostel: req.Hostel,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	if !user.CheckPassword(password) {
		return nil, models.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) GetAll(page, pageSize int) ([]models.User, int64, error) {
	return s.repo.FindAll(page, pageSize)
}

func (s *UserService) Update(user *models.User) error {
	return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
