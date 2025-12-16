package repository

import (
	"darulabror/internal/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	//Manage Admins by Superadmin
	CreateAdmin(admin models.Admin) error
	GetAllAdmins(page, limit int) ([]models.Admin, int64, error)
	GetAdminByID(id uint) (models.Admin, error)
	GetAdminByEmail(email string) (models.Admin, error)
	UpdateAdmin(admin models.Admin) error
	DeleteAdmin(id uint) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) CreateAdmin(admin models.Admin) error {
	return r.db.Create(&admin).Error
}

func (r *adminRepository) GetAllAdmins(page, limit int) ([]models.Admin, int64, error) {
	var (
		admins []models.Admin
		total  int64
	)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err := r.db.Model(&models.Admin{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Order("id DESC").Limit(limit).Offset(offset).Find(&admins).Error
	return admins, total, err
}

func (r *adminRepository) GetAdminByID(id uint) (models.Admin, error) {
	var admin models.Admin
	err := r.db.First(&admin, id).Error
	return admin, err
}

func (r *adminRepository) GetAdminByEmail(email string) (models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	return admin, err
}

func (r *adminRepository) UpdateAdmin(admin models.Admin) error {
	return r.db.Save(&admin).Error
}

func (r *adminRepository) DeleteAdmin(id uint) error {
	return r.db.Delete(&models.Admin{}, id).Error
}
