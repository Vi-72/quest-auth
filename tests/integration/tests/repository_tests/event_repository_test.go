// REPOSITORY LAYER INTEGRATION TESTS
// Tests for repository implementations and database interactions

//go:build integration

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Vi-72/quest-auth/internal/adapters/out/postgres/eventrepo"
	"github.com/Vi-72/quest-auth/internal/core/domain/model/auth"
)

func (s *Suite) TestEventRepository_Publish_SavesEvent() {
	ctx := context.Background()

	// Pre-condition: repository
	repo := eventrepo.NewRepository(s.TestDIContainer.DB)

	userID := uuid.New()
	ev := auth.NewUserRegistered(userID, "repo-event@example.com", "+1234567890", time.Now())

	// Act: publish sync
	err := repo.Publish(ctx, ev)
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
