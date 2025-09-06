// REPOSITORY LAYER INTEGRATION TESTS
// Tests for repository implementations and database interactions

//go:build integration

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"quest-auth/internal/adapters/out/postgres/eventrepo"
	"quest-auth/internal/core/domain/model/auth"
	"quest-auth/internal/core/ports"
)

func (s *Suite) TestEventRepository_Publish_SavesEvent() {
	ctx := context.Background()

	// Pre-condition: repository
	tracker := s.TestDIContainer.UnitOfWork.(ports.Tracker)
	repo, err := eventrepo.NewRepository(tracker, 5)
	s.Require().NoError(err)

	userID := uuid.New()
	ev := auth.NewUserRegistered(userID, "repo-event@example.com", "+1234567890", time.Now())

	// Act: publish sync
	err = repo.Publish(ctx, ev)
	s.Require().NoError(err)

	// Assert: stored in DB
	var rows []eventrepo.EventDTO
	q := s.TestDIContainer.DB.Where("aggregate_id = ?", userID.String())
	s.Require().NoError(q.Find(&rows).Error)
	s.Assert().GreaterOrEqual(len(rows), 1)
	var types []string
	for _, r := range rows {
		types = append(types, r.EventType)
	}
	s.Contains(types, "user.registered")
}

func (s *Suite) TestEventRepository_PublishAsync_SavesEvent() {
	ctx := context.Background()

	// Pre-condition: repository
	tracker := s.TestDIContainer.UnitOfWork.(ports.Tracker)
	repo, err := eventrepo.NewRepository(tracker, 1)
	s.Require().NoError(err)

	userID := uuid.New()
	ev := auth.NewUserLoggedIn(userID, time.Now())

	// Act: publish async
	repo.PublishAsync(ctx, ev)

	// Wait for async write using polling (up to 2s)
	deadline := time.Now().Add(2 * time.Second)
	for {
		var rows []eventrepo.EventDTO
		q := s.TestDIContainer.DB.Where("aggregate_id = ?", userID.String())
		if err := q.Find(&rows).Error; err == nil && len(rows) >= 1 {
			break
		}
		if time.Now().After(deadline) {
			s.Require().Fail("timeout waiting for async event")
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Assert: stored in DB with expected type
	var rows []eventrepo.EventDTO
	q := s.TestDIContainer.DB.Where("aggregate_id = ?", userID.String())
	s.Require().NoError(q.Find(&rows).Error)
	s.Assert().GreaterOrEqual(len(rows), 1)
	var types2 []string
	for _, r := range rows {
		types2 = append(types2, r.EventType)
	}
	s.Contains(types2, "user.login")
}
