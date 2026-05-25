package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"artale_market/model"
	"artale_market/repository"
)

type MemberService interface {
	Authenticate(username, password string) (*model.Member, error)
	Register(nickname, username, password, email string) (*model.Member, error)
	GetMe(id uint) (*model.Member, error)
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

func md5Hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func (s *memberService) Authenticate(username, password string) (*model.Member, error) {
	member, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("帳號或密碼錯誤")
	}
	if member.Status == 0 {
		return nil, errors.New("帳號已停用")
	}
	if member.Password != md5Hash(password) {
		return nil, errors.New("帳號或密碼錯誤")
	}
	return member, nil
}

func (s *memberService) Register(nickname, username, password, email string) (*model.Member, error) {
	if _, err := s.repo.FindByUsername(username); err == nil {
		return nil, errors.New("帳號已被使用")
	}
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, errors.New("信箱已被使用")
	}

	member := &model.Member{
		Nickname: nickname,
		Username: username,
		Password: md5Hash(password),
		Email:    email,
		Status:   1,
	}
	if err := s.repo.Create(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *memberService) GetMe(id uint) (*model.Member, error) {
	member, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
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
