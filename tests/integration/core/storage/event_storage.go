package storage

import (
	"context"
	"gorm.io/gorm"

	"quest-auth/internal/adapters/out/postgres/eventrepo"
)

// EventStorage provides helpers to work with events in tests
type EventStorage struct {
	db *gorm.DB
}

func NewEventStorage(db *gorm.DB) *EventStorage { return &EventStorage{db: db} }

func (s *EventStorage) GetEventsByType(ctx context.Context, eventType string) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	return events, s.db.WithContext(ctx).Where("event_type = ?", eventType).Order("created_at ASC").Find(&events).Error
}
