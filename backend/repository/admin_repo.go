package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type AdminRepository interface {
	FindByUsername(username string) (*model.AdminUser, error)
	FindAll() ([]model.AdminUser, error)
	FindByID(id uint) (*model.AdminUser, error)
	Create(admin *model.AdminUser) error
	Update(admin *model.AdminUser) error
	Delete(id uint) error
	Count() (int64, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) FindByUsername(username string) (*model.AdminUser, error) {
	var admin model.AdminUser
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) FindAll() ([]model.AdminUser, error) {
	var admins []model.AdminUser
	if err := r.db.Order("id asc").Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *adminRepository) FindByID(id uint) (*model.AdminUser, error) {
	var admin model.AdminUser
	if err := r.db.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) Create(admin *model.AdminUser) error {
	return r.db.Create(admin).Error
}

func (r *adminRepository) Update(admin *model.AdminUser) error {
	return r.db.Save(admin).Error
}

func (r *adminRepository) Delete(id uint) error {
	return r.db.Delete(&model.AdminUser{}, id).Error
}

func (r *adminRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&model.AdminUser{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
