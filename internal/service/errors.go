package service

import "errors"

var (
	// Admin service errors
	ErrNotFoundAdmin = errors.New("admin not found")
	ErrInvalidAdmin  = errors.New("invalid admin")
	ErrCreateAdmin   = errors.New("failed to create admin")
	// Article service errors
	ErrNotFoundArticle = errors.New("article not found")
	ErrCreateArticle   = errors.New("failed to create article")
	ErrUpdateArticle   = errors.New("failed to update article")
	// Registration service errors public
	ErrCreateRegistration = errors.New("failed to create registration")
	// Registration service errors admin
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAdminInactive      = errors.New("admin is inactive")
	// Registration service errors additional
	ErrRegistrationEmailExists = errors.New("registration email already used")
	ErrRegistrationNISNExists  = errors.New("registration nisn already used")
)
