package eventrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Vi-72/quest-auth/internal/core/domain/model/auth"
	"github.com/Vi-72/quest-auth/internal/core/ports"
	"github.com/Vi-72/quest-auth/internal/pkg/ddd"
	"github.com/Vi-72/quest-auth/internal/pkg/errs"
)

var _ ports.EventPublisher = &Repository{}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	var dtos []EventDTO
	for _, event := range events {
		dto, err := r.domainEventToDTO(event)
		if err != nil {
			return errs.WrapInfrastructureError("failed to convert event to DTO", err)
		}
		dtos = append(dtos, dto)
	}

	db := r.db.WithContext(ctx)
	for i := range dtos {
		if err := db.Create(&dtos[i]).Error; err != nil {
			return errs.WrapInfrastructureError("failed to save event", err)
		}
	}

	return nil
}

func (r *Repository) domainEventToDTO(event ddd.DomainEvent) (EventDTO, error) {
	dto := EventDTO{
		ID:        event.GetID().String(),
		EventType: event.GetName(),
		CreatedAt: time.Now(),
	}

	switch e := event.(type) {
	case auth.UserRegistered,
		auth.UserPhoneChanged,
		auth.UserNameChanged,
		auth.UserPasswordChanged,
		auth.UserLoggedIn:
		agg, ok := e.(interface {
			GetAggregateID() uuid.UUID
		})
		if !ok {
			return EventDTO{}, errs.NewDomainValidationError("eventSerialization", "event missing AggregateID")
		}

		dto.AggregateID = agg.GetAggregateID().String()

		data, err := MarshalEventData(e)
		if err != nil {
			return EventDTO{}, err
		}
		dto.Data = data

	default:
		if agg, ok := e.(interface {
			GetAggregateID() uuid.UUID
		}); ok {
			dto.AggregateID = agg.GetAggregateID().String()
		} else {
			dto.AggregateID = event.GetID().String()
		}

		data, err := MarshalEventData(event)
		if err != nil {
			return EventDTO{}, errs.NewDomainValidationError("eventSerialization", "failed to serialize unknown event type")
		}
		dto.Data = data
	}

	return dto, nil
}
