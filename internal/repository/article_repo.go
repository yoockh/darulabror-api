package repository

import (
	"darulabror/internal/models"
	"darulabror/internal/utils"

	"gorm.io/gorm"
)

type ArticleRepo interface {
	Create(article models.Article) error
	GetAll(page, limit int) ([]models.Article, int64, error)
	GetPublished(page, limit int) ([]models.Article, int64, error)
	GetByID(id uint) (models.Article, error)
	Update(article models.Article) error
	Delete(id uint) error
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) ArticleRepo {
	return &articleRepo{db: db}
}

func (a *articleRepo) Create(article models.Article) error {
	return a.db.Create(&article).Error
}

func (a *articleRepo) GetAll(page, limit int) ([]models.Article, int64, error) {
	var (
		articles []models.Article
		total    int64
	)

	_, limit, offset := utils.NormalizePageLimit(page, limit)

	if err := a.db.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := a.db.Order("id DESC").Limit(limit).Offset(offset).Find(&articles).Error
	return articles, total, err
}

func (a *articleRepo) GetPublished(page, limit int) ([]models.Article, int64, error) {
	var (
		articles []models.Article
		total    int64
	)

	_, limit, offset := utils.NormalizePageLimit(page, limit)

	q := a.db.Model(&models.Article{}).Where("status = ?", "published")

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order("id DESC").Limit(limit).Offset(offset).Find(&articles).Error
	return articles, total, err
}

func (a *articleRepo) GetByID(id uint) (models.Article, error) {
	var article models.Article
	err := a.db.First(&article, id).Error
	return article, err
}

func (a *articleRepo) Update(article models.Article) error {
	return a.db.Save(&article).Error
}

func (a *articleRepo) Delete(id uint) error {
	return a.db.Delete(&models.Article{}, id).Error
}
