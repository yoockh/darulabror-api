package dto

import (
	"darulabror/internal/models"

	"gorm.io/datatypes"
)

type ArticleDTO struct {
	ID      uint           `json:"id" validate:"omitempty"`
	Title   string         `json:"title" validate:"required,min=3,max=100"`
	Content datatypes.JSON `json:"content" validate:"required"`
	Author  string         `json:"author" validate:"required,min=3,max=50"`
	Status  string         `json:"status" validate:"omitempty,oneof=draft published"`

	CreatedAt int64 `json:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty"`
}

func ArticleDTOToModel(dto ArticleDTO) (models.Article, error) {
	return models.Article{
		Title:   dto.Title,
		Content: dto.Content,
		Author:  dto.Author,
		Status:  dto.Status,
	}, nil
}

func ArticleModelToDTO(article models.Article) ArticleDTO {
	return ArticleDTO{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Author:    article.Author,
		Status:    article.Status,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
}
