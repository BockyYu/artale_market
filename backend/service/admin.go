package service

import (
	"errors"
	"regexp"

	"artale_market/model"
	"artale_market/repository"

	"golang.org/x/crypto/bcrypt"
)

var alphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

type AdminService interface {
	Authenticate(username, password string) (*model.AdminUser, error)
	GetByID(id uint) (*model.AdminUser, error)
	ListAdmins() ([]model.AdminUser, error)
	CreateAdmin(username, password, role string) (*model.AdminUser, error)
	UpdateAdmin(id uint, username, password, role string) (*model.AdminUser, error)
	DeleteAdmin(id uint) error
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{repo: repo}
}

func (s *adminService) GetByID(id uint) (*model.AdminUser, error) {
	return s.repo.FindByID(id)
}

func (s *adminService) Authenticate(username, password string) (*model.AdminUser, error) {
	admin, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return admin, nil
}

func (s *adminService) ListAdmins() ([]model.AdminUser, error) {
	return s.repo.FindAll()
}

func (s *adminService) CreateAdmin(username, password, role string) (*model.AdminUser, error) {
	if !alphanumeric.MatchString(username) || !alphanumeric.MatchString(password) {
		return nil, errors.New("username and password must be alphanumeric")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if role == "" {
		role = "admin"
	}
	admin := &model.AdminUser{
		Username: username,
		Password: string(hash),
		Role:     role,
	}
	if err := s.repo.Create(admin); err != nil {
		return nil, errors.New("username already exists")
	}
	return admin, nil
}

func (s *adminService) UpdateAdmin(id uint, username, password, role string) (*model.AdminUser, error) {
	admin, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("admin not found")
	}
	if username != "" {
		if !alphanumeric.MatchString(username) {
			return nil, errors.New("username must be alphanumeric")
		}
		admin.Username = username
	}
	if password != "" {
		if !alphanumeric.MatchString(password) {
			return nil, errors.New("password must be alphanumeric")
		}
		if len(password) < 6 {
			return nil, errors.New("password must be at least 6 characters")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		admin.Password = string(hash)
	}
	if role != "" {
		admin.Role = role
	}
	if err := s.repo.Update(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *adminService) DeleteAdmin(id uint) error {
	count, err := s.repo.Count()
	if err != nil {
		return err
	}
	if count <= 1 {
		return errors.New("cannot delete the last admin")
	}
	return s.repo.Delete(id)
}
