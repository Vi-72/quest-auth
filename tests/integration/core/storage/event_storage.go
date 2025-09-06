package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quest-auth/internal/adapters/out/postgres/eventrepo"
)

// EventStorage provides helpers to work with events in tests
type EventStorage struct {
	db *gorm.DB
}

func NewEventStorage(db *gorm.DB) *EventStorage { return &EventStorage{db: db} }

func (s *EventStorage) GetEventByID(ctx context.Context, id string) (*eventrepo.EventDTO, error) {
	var e eventrepo.EventDTO
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *EventStorage) GetAllEvents(ctx context.Context) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	return events, s.db.WithContext(ctx).Order("created_at ASC").Find(&events).Error
}

func (s *EventStorage) GetEventsByType(ctx context.Context, eventType string) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	return events, s.db.WithContext(ctx).Where("event_type = ?", eventType).Order("created_at ASC").Find(&events).Error
}

func (s *EventStorage) GetEventsByAggregateID(ctx context.Context, aggregateID uuid.UUID) ([]eventrepo.EventDTO, error) {
	var events []eventrepo.EventDTO
	return events, s.db.WithContext(ctx).Where("aggregate_id = ?", aggregateID.String()).Order("created_at ASC").Find(&events).Error
}

func (s *EventStorage) CountEvents(ctx context.Context) (int64, error) {
	var count int64
	return count, s.db.WithContext(ctx).Model(&eventrepo.EventDTO{}).Count(&count).Error
}

func (s *EventStorage) CountEventsByType(ctx context.Context, eventType string) (int64, error) {
	var count int64
	return count, s.db.WithContext(ctx).Model(&eventrepo.EventDTO{}).Where("event_type = ?", eventType).Count(&count).Error
}

func (s *EventStorage) WaitForEventsOfType(ctx context.Context, eventType string, expectedCount int64, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		count, err := s.CountEventsByType(ctx, eventType)
		if err != nil {
			return err
		}
		if count >= expectedCount {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return gorm.ErrRecordNotFound
}
