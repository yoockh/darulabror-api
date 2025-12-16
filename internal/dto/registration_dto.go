package dto

import (
	"darulabror/internal/models"
	"time"
)

type RegistrationDTO struct {
	ID uint `json:"id" validate:"omitempty"`

	StudentType models.StudentType `json:"student_type" validate:"required,oneof=new transfer"`
	FullName    string             `json:"full_name" validate:"required,min=3,max=100"`
	Email       string             `json:"email" validate:"required,email"`
	Phone       string             `json:"phone" validate:"required,min=10,max=13"`

	Gender       models.Gender `json:"gender" validate:"required,oneof=male female"`
	PlaceOfBirth string        `json:"place_of_birth" validate:"required,min=3,max=100"`
	DateOfBirth  string        `json:"date_of_birth" validate:"required,datetime=2006-01-02"`

	Address      string `json:"address" validate:"required,min=3,max=255"`
	OriginSchool string `json:"origin_school" validate:"required,min=3,max=100"`
	NISN         string `json:"nisn" validate:"required,len=10"`

	FatherName        string `json:"father_name" validate:"required,min=3,max=100"`
	FatherOccupation  string `json:"father_occupation" validate:"required,min=3,max=100"`
	PhoneFather       string `json:"phone_father" validate:"required,min=10,max=13"`
	DateOfBirthFather string `json:"date_of_birth_father" validate:"required,datetime=2006-01-02"`

	MotherName        string `json:"mother_name" validate:"required,min=3,max=100"`
	MotherOccupation  string `json:"mother_occupation" validate:"required,min=3,max=100"`
	PhoneMother       string `json:"phone_mother" validate:"required,min=10,max=13"`
	DateOfBirthMother string `json:"date_of_birth_mother" validate:"required,datetime=2006-01-02"`

	CreatedAt string `json:"created_at,omitempty"`
}

const dateLayout = "2006-01-02"

func RegistrationDTOToModel(d RegistrationDTO) (models.Registration, error) {
	dob, err := time.Parse(dateLayout, d.DateOfBirth)
	if err != nil {
		return models.Registration{}, err
	}
	dobF, err := time.Parse(dateLayout, d.DateOfBirthFather)
	if err != nil {
		return models.Registration{}, err
	}
	dobM, err := time.Parse(dateLayout, d.DateOfBirthMother)
	if err != nil {
		return models.Registration{}, err
	}

	return models.Registration{
		ID:                d.ID,
		StudentType:       d.StudentType,
		Gender:            d.Gender,
		Email:             d.Email,
		FullName:          d.FullName,
		Phone:             d.Phone,
		PlaceOfBirth:      d.PlaceOfBirth,
		DateOfBirth:       dob,
		Address:           d.Address,
		OriginSchool:      d.OriginSchool,
		NISN:              d.NISN,
		FatherName:        d.FatherName,
		FatherOccupation:  d.FatherOccupation,
		PhoneFather:       d.PhoneFather,
		DateOfBirthFather: dobF,
		MotherName:        d.MotherName,
		MotherOccupation:  d.MotherOccupation,
		PhoneMother:       d.PhoneMother,
		DateOfBirthMother: dobM,
	}, nil
}

func RegistrationModelToDTO(m models.Registration) RegistrationDTO {
	return RegistrationDTO{
		ID:                m.ID,
		StudentType:       m.StudentType,
		FullName:          m.FullName,
		Email:             m.Email,
		Phone:             m.Phone,
		Gender:            m.Gender,
		PlaceOfBirth:      m.PlaceOfBirth,
		DateOfBirth:       m.DateOfBirth.Format(dateLayout),
		Address:           m.Address,
		OriginSchool:      m.OriginSchool,
		NISN:              m.NISN,
		FatherName:        m.FatherName,
		FatherOccupation:  m.FatherOccupation,
		PhoneFather:       m.PhoneFather,
		DateOfBirthFather: m.DateOfBirthFather.Format(dateLayout),
		MotherName:        m.MotherName,
		MotherOccupation:  m.MotherOccupation,
		PhoneMother:       m.PhoneMother,
		DateOfBirthMother: m.DateOfBirthMother.Format(dateLayout),
		CreatedAt:         m.CreatedAt.Format(time.RFC3339),
	}
}
