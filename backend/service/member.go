package service

import (
	"errors"

	"artale_market/model"
	"artale_market/repository"

	"golang.org/x/crypto/bcrypt"
)

type MemberService interface {
	Authenticate(username, password string) (*model.Member, error)
	ListMembers(page, pageSize int, search string) ([]model.Member, int64, error)
	UpdateStatus(id uint, status int) (*model.Member, error)
	Delete(id uint) error
}

type memberService struct {
	repo repository.MemberRepository
}

func NewMemberService(repo repository.MemberRepository) MemberService {
	return &memberService{repo: repo}
}

func (s *memberService) Authenticate(username, password string) (*model.Member, error) {
	member, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("帳號或密碼錯誤")
	}
	if member.Status == 0 {
		return nil, errors.New("帳號已停用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password)); err != nil {
		return nil, errors.New("帳號或密碼錯誤")
	}
	return member, nil
}

func (s *memberService) ListMembers(page, pageSize int, search string) ([]model.Member, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.FindAll(page, pageSize, search)
}

func (s *memberService) UpdateStatus(id uint, status int) (*model.Member, error) {
	if status != 0 && status != 1 {
		return nil, errors.New("status must be 0 or 1")
	}
	member, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("member not found")
	}
	member.Status = status
	if err := s.repo.Update(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *memberService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("member not found")
	}
	return s.repo.Delete(id)
}
