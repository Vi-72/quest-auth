package userrepo

import (
	"errors"

	"github.com/Vi-72/quest-auth/internal/core/domain/model/auth"
	"github.com/Vi-72/quest-auth/internal/core/domain/model/kernel"
	"github.com/Vi-72/quest-auth/internal/core/ports"
	"github.com/Vi-72/quest-auth/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create сохраняет нового пользователя
func (r *Repository) Create(user *auth.User) error {
	dto := FromEntity(user)

	if err := r.db.Create(&dto).Error; err != nil {
		return errs.WrapInfrastructureError("creating user", err)
	}

	return nil
}

// GetByID находит пользователя по ID
func (r *Repository) GetByID(id uuid.UUID) (*auth.User, error) {
	var dto UserDTO
	err := r.db.Where("id = ?", id).First(&dto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("user", id.String())
		}
		return nil, errs.WrapInfrastructureError("getting user by id", err)
	}

	return dto.ToEntity()
}

// GetByEmail находит пользователя по email
func (r *Repository) GetByEmail(email kernel.Email) (*auth.User, error) {
	var dto UserDTO
	err := r.db.Where("email = ?", email.String()).First(&dto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("user", email.String())
		}
		return nil, errs.WrapInfrastructureError("getting user by email", err)
	}

	return dto.ToEntity()
}

// GetByPhone находит пользователя по телефону
func (r *Repository) GetByPhone(phone kernel.Phone) (*auth.User, error) {
	var dto UserDTO
	err := r.db.Where("phone = ?", phone.String()).First(&dto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("user", phone.String())
		}
		return nil, errs.WrapInfrastructureError("getting user by phone", err)
	}

	return dto.ToEntity()
}

// Update обновляет существующего пользователя
func (r *Repository) Update(user *auth.User) error {
	dto := FromEntity(user)

	result := r.db.Where("id = ?", user.ID()).Updates(&dto)
	if result.Error != nil {
		return errs.WrapInfrastructureError("updating user", result.Error)
	}

	if result.RowsAffected == 0 {
		return errs.NewNotFoundError("user", user.ID().String())
	}

	return nil
}

// Delete удаляет пользователя
func (r *Repository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&UserDTO{})

	if result.Error != nil {
		return errs.WrapInfrastructureError("deleting user", result.Error)
	}

	if result.RowsAffected == 0 {
		return errs.NewNotFoundError("user", id.String())
	}

	return nil
}

// EmailExists проверяет существование email
func (r *Repository) EmailExists(email kernel.Email) (bool, error) {
	var count int64
	err := r.db.Model(&UserDTO{}).Where("email = ?", email.String()).Count(&count).Error

	if err != nil {
		return false, errs.WrapInfrastructureError("checking email existence", err)
	}

	return count > 0, nil
}

// PhoneExists проверяет существование телефона
func (r *Repository) PhoneExists(phone kernel.Phone) (bool, error) {
	var count int64
	err := r.db.Model(&UserDTO{}).Where("phone = ?", phone.String()).Count(&count).Error

	if err != nil {
		return false, errs.WrapInfrastructureError("checking phone existence", err)
	}

	return count > 0, nil
}

// Compile-time check that Repository implements UserRepository
var _ ports.UserRepository = (*Repository)(nil)
