package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type MemberRepository interface {
	FindByUsername(username string) (*model.Member, error)
	FindByEmail(email string) (*model.Member, error)
	FindAll(page, pageSize int, search string) ([]model.Member, int64, error)
	FindByID(id uint) (*model.Member, error)
	Create(member *model.Member) error
	Update(member *model.Member) error
	Delete(id uint) error
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) FindByUsername(username string) (*model.Member, error) {
	var member model.Member
	if err := r.db.Where("username = ?", username).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) FindByEmail(email string) (*model.Member, error) {
	var member model.Member
	if err := r.db.Where("email = ?", email).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) FindAll(page, pageSize int, search string) ([]model.Member, int64, error) {
	var members []model.Member
	var total int64

	q := r.db.Model(&model.Member{})
	if search != "" {
		q = q.Where("username ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := q.Order("id asc").Offset(offset).Limit(pageSize).Find(&members).Error; err != nil {
		return nil, 0, err
	}
	return members, total, nil
}

func (r *memberRepository) FindByID(id uint) (*model.Member, error) {
	var member model.Member
	if err := r.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) Create(member *model.Member) error {
	return r.db.Create(member).Error
}

func (r *memberRepository) Update(member *model.Member) error {
	return r.db.Save(member).Error
}

func (r *memberRepository) Delete(id uint) error {
	return r.db.Delete(&model.Member{}, id).Error
}
