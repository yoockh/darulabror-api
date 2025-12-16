package models

import "time"

type StudentType string

const (
	StudentNew      StudentType = "new"
	StudentTransfer StudentType = "transfer"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Registration struct {
	ID uint `gorm:"primaryKey" json:"id"`

	StudentType StudentType `gorm:"type:text;check:student_type IN ('new','transfer');not null" json:"student_type"`
	Gender      Gender      `gorm:"type:text;check:gender IN ('male','female');not null" json:"gender"`

	Email    string `gorm:"not null;uniqueIndex" json:"email"`
	FullName string `gorm:"not null" json:"full_name"`
	Phone    string `gorm:"not null" json:"phone"`

	PlaceOfBirth string    `gorm:"not null" json:"place_of_birth"`
	DateOfBirth  time.Time `gorm:"not null" json:"date_of_birth"`

	Address      string `gorm:"type:text;not null" json:"address"`
	OriginSchool string `gorm:"not null" json:"origin_school"`

	NISN string `gorm:"not null;uniqueIndex" json:"nisn"`

	FatherName        string    `gorm:"not null" json:"father_name"`
	FatherOccupation  string    `gorm:"not null" json:"father_occupation"`
	PhoneFather       string    `gorm:"not null" json:"phone_father"`
	DateOfBirthFather time.Time `gorm:"not null" json:"date_of_birth_father"`

	MotherName        string    `gorm:"not null" json:"mother_name"`
	MotherOccupation  string    `gorm:"not null" json:"mother_occupation"`
	PhoneMother       string    `gorm:"not null" json:"phone_mother"`
	DateOfBirthMother time.Time `gorm:"not null" json:"date_of_birth_mother"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
