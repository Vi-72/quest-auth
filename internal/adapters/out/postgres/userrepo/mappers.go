package userrepo

import (
	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/domain/model/kernel"
	"quest-auth/internal/pkg/ddd"
	"quest-auth/internal/pkg/errs"
)

// ToEntity преобразует DTO в доменную сущность User
func (dto UserDTO) ToEntity() (auth.User, error) {
	email, err := kernel.NewEmail(dto.Email)
	if err != nil {
		return auth.User{}, errs.WrapInfrastructureError("mapping user email", err)
	}

	phone, err := kernel.NewPhone(dto.Phone)
	if err != nil {
		return auth.User{}, errs.WrapInfrastructureError("mapping user phone", err)
	}

	user := auth.User{
		BaseAggregate: ddd.NewBaseAggregate(dto.ID),
		Email:         email,
		Phone:         phone,
		Name:          dto.Name,
		PasswordHash:  dto.PasswordHash,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}

	return user, nil
}

// FromEntity преобразует доменную сущность User в DTO
func FromEntity(user auth.User) UserDTO {
	return UserDTO{
		ID:           user.ID(),
		Email:        user.Email.String(),
		Phone:        user.Phone.String(),
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
