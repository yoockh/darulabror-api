package service

import (
	"darulabror/internal/dto"
	"darulabror/internal/repository"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RegistrationService interface {
	// Public
	CreateRegistration(regDTO dto.RegistrationDTO) error

	// Admin
	GetAllRegistrations(page, limit int) ([]dto.RegistrationDTO, int64, error)
	GetRegistrationByID(id uint) (dto.RegistrationDTO, error)
	DeleteRegistration(id uint) error
}

type registrationService struct {
	repo repository.RegistrationRepo
}

func NewRegistrationService(repo repository.RegistrationRepo) RegistrationService {
	return &registrationService{repo: repo}
}

func (s *registrationService) CreateRegistration(regDTO dto.RegistrationDTO) error {
	// uniqueness checks
	existsEmail, err := s.repo.ExistsByEmail(regDTO.Email)
	if err != nil {
		logrus.WithError(err).WithField("email", regDTO.Email).Error("failed check registration email")
		return err
	}
	if existsEmail {
		return ErrRegistrationEmailExists
	}

	existsNISN, err := s.repo.ExistsByNISN(regDTO.NISN)
	if err != nil {
		logrus.WithError(err).WithField("nisn", regDTO.NISN).Error("failed check registration nisn")
		return err
	}
	if existsNISN {
		return ErrRegistrationNISNExists
	}

	reg, err := dto.RegistrationDTOToModel(regDTO)
	if err != nil {
		logrus.WithError(err).Error("failed convert RegistrationDTO to model")
		return err
	}

	if err := s.repo.Create(reg); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"email": reg.Email,
			"nisn":  reg.NISN,
		}).Error("failed create registration")
		return ErrCreateRegistration
	}

	logrus.WithFields(logrus.Fields{
		"email": reg.Email,
		"nisn":  reg.NISN,
	}).Info("registration created")
	return nil
}

func (s *registrationService) GetAllRegistrations(page, limit int) ([]dto.RegistrationDTO, int64, error) {
	regs, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		logrus.WithError(err).Error("failed get all registrations")
		return nil, 0, err
	}

	out := make([]dto.RegistrationDTO, 0, len(regs))
	for _, r := range regs {
		out = append(out, dto.RegistrationModelToDTO(r))
	}
	return out, total, nil
}

func (s *registrationService) GetRegistrationByID(id uint) (dto.RegistrationDTO, error) {
	reg, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.RegistrationDTO{}, errors.New("registration not found")
		}
		logrus.WithError(err).WithField("id", id).Error("failed get registration by id")
		return dto.RegistrationDTO{}, err
	}
	return dto.RegistrationModelToDTO(reg), nil
}

func (s *registrationService) DeleteRegistration(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed delete registration")
		return err
	}
	logrus.WithField("id", id).Info("registration deleted")
	return nil
}
