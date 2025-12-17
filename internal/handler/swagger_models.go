package handler

import "darulabror/internal/dto"

type PaginationMeta struct {
	Page  int   `json:"page" example:"1"`
	Limit int   `json:"limit" example:"10"`
	Total int64 `json:"total" example:"123"`
}

type ListResponseData[T any] struct {
	Items []T            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type SuccessResponse[T any] struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"OK"`
	Data    T      `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"something went wrong"`
}

// ===== Requests (named types so Swagger shows templates) =====

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"admin@darulabror.com"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"StrongPassword123"`
}

type AdminLoginResponseData struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Admin dto.AdminDTO `json:"admin"`
}

type ContactCreateRequest struct {
	Email   string `json:"email" validate:"required,email" example:"user@example.com"`
	Subject string `json:"subject" validate:"required,min=3,max=150" example:"Question about registration"`
	Message string `json:"message" validate:"required,min=3,max=2000" example:"Hello, I would like to ask..."`
}

type ContactUpdateRequest = ContactCreateRequest

// ===== Typed responses used in annotations =====

type ArticleListResponse = SuccessResponse[ListResponseData[dto.ArticleDTO]]
type RegistrationListResponse = SuccessResponse[ListResponseData[dto.RegistrationDTO]]
type AdminListResponse = SuccessResponse[ListResponseData[dto.AdminDTO]]

type ContactListItem struct {
	ID        uint   `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	Subject   string `json:"subject" example:"Question"`
	Message   string `json:"message" example:"Hello..."`
	CreatedAt int64  `json:"created_at" example:"1734567890"`
}

type ContactListResponse = SuccessResponse[ListResponseData[ContactListItem]]

type AdminLoginResponse = SuccessResponse[AdminLoginResponseData]
