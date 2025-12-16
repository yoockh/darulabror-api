package service

import (
	"context"
	"darulabror/internal/dto"
	"darulabror/internal/repository"
	"errors"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ArticleService interface {
	// Public
	GetPublishedArticles(page, limit int) ([]dto.ArticleDTO, int64, error)
	GetPublishedArticleByID(id uint) (dto.ArticleDTO, error)

	// Admin
	CreateArticle(articleDTO dto.ArticleDTO) error
	GetAllArticles(page, limit int) ([]dto.ArticleDTO, int64, error)
	UpdateArticle(id uint, articleDTO dto.ArticleDTO) error
	DeleteArticle(id uint) error

	// ======================
	//  METHODS FOR GCS
	// ======================
	UploadArticleMedia(ctx context.Context, file io.Reader, fileName string) (string, error)
	GetArticleMediaURL(ctx context.Context, objectName string) (string, error)
}

type articleService struct {
	repo         repository.ArticleRepo
	privateStore repository.GCPStorageRepo
}

func NewArticleService(repo repository.ArticleRepo, privateStore repository.GCPStorageRepo) ArticleService {
	return &articleService{
		repo:         repo,
		privateStore: privateStore,
	}
}

func (s *articleService) CreateArticle(articleDTO dto.ArticleDTO) error {
	if articleDTO.Status == "" {
		articleDTO.Status = "draft"
	}

	article, err := dto.ArticleDTOToModel(articleDTO)
	if err != nil {
		logrus.WithError(err).Error("failed convert ArticleDTO to model")
		return err
	}

	if err := s.repo.Create(article); err != nil {
		logrus.WithError(err).WithField("title", article.Title).Error("failed to create article")
		return ErrCreateArticle
	}

	logrus.WithField("title", article.Title).Info("article created")
	return nil
}

func (s *articleService) GetAllArticles(page, limit int) ([]dto.ArticleDTO, int64, error) {
	articles, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		logrus.WithError(err).Error("failed get all articles")
		return nil, 0, err
	}

	out := make([]dto.ArticleDTO, 0, len(articles))
	for _, a := range articles {
		out = append(out, dto.ArticleModelToDTO(a))
	}
	return out, total, nil
}

func (s *articleService) GetPublishedArticles(page, limit int) ([]dto.ArticleDTO, int64, error) {
	articles, total, err := s.repo.GetPublished(page, limit)
	if err != nil {
		logrus.WithError(err).Error("failed get published articles")
		return nil, 0, err
	}

	out := make([]dto.ArticleDTO, 0, len(articles))
	for _, a := range articles {
		out = append(out, dto.ArticleModelToDTO(a))
	}
	return out, total, nil
}

func (s *articleService) GetPublishedArticleByID(id uint) (dto.ArticleDTO, error) {
	article, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ArticleDTO{}, ErrNotFoundArticle
		}
		logrus.WithError(err).WithField("id", id).Error("failed get article by id")
		return dto.ArticleDTO{}, err
	}

	if article.Status != "published" {
		return dto.ArticleDTO{}, ErrNotFoundArticle
	}
	return dto.ArticleModelToDTO(article), nil
}

func (s *articleService) UpdateArticle(id uint, articleDTO dto.ArticleDTO) error {
	article, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFoundArticle
		}
		return err
	}

	article.Title = articleDTO.Title
	article.Content = articleDTO.Content
	article.Author = articleDTO.Author
	if articleDTO.Status != "" {
		article.Status = articleDTO.Status
	}

	if err := s.repo.Update(article); err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed update article")
		return ErrUpdateArticle
	}

	logrus.WithField("id", id).Info("article updated")
	return nil
}

func (s *articleService) DeleteArticle(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed delete article")
		return err
	}
	logrus.WithField("id", id).Info("article deleted")
	return nil
}

// ======================
//  METHODS FOR GCS
// ======================

func (s *articleService) UploadArticleMedia(ctx context.Context, file io.Reader, fileName string) (string, error) {
	objectName, err := s.privateStore.UploadFile(ctx, file, fileName)
	if err != nil {
		logrus.WithError(err).WithField("fileName", fileName).Error("failed upload article media")
		return "", err
	}
	return objectName, nil
}

func (s *articleService) GetArticleMediaURL(ctx context.Context, objectName string) (string, error) {
	url, err := s.privateStore.GenerateSignedURL(ctx, objectName, 10*time.Minute)
	if err != nil {
		logrus.WithError(err).WithField("objectName", objectName).Error("failed generate article signed url")
		return "", err
	}
	return url, nil
}
